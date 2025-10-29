// ---------------------------
// pkg/runner/runner.go
// ---------------------------
package runner

import (
	"encoding/json"
	// "errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"os/exec"

	"infractest/internal/terraform"
	"infractest/pkg/assert"
	"infractest/pkg/mocks"
	"infractest/pkg/parser"
)

// TestCaseResult captures a single test's results.
type TestCaseResult struct {
	File       string                    `json:"file"`
	Test       string                    `json:"test"`
	Assertions []assert.Result          `json:"assertions"`
	Passed     bool                      `json:"passed"`
	Logs       string                    `json:"logs,omitempty"`
}

// RunDirectory parses and runs tests in a directory.
func RunDirectory(dir string, mode string) ([]TestCaseResult, error) {
	parsed, err := parser.ParseDirectory(dir)
	if err != nil { return nil, err }

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []TestCaseResult{}
	errorChan := make(chan error, 1)

	for path, tf := range parsed {
		for _, t := range tf.Tests {
			wg.Add(1)
			go func(p string, tb parser.TestBlock) {
				defer wg.Done()
				res := TestCaseResult{File: filepath.Base(p), Test: tb.Name}

				// 1) Create sandbox dir
				tmpDir, err := ioutil.TempDir("", "terraspec-")
				if err != nil {
					res.Passed = false
					res.Logs = err.Error()
					mu.Lock(); results = append(results, res); mu.Unlock()
					return
				}
				defer os.RemoveAll(tmpDir)

				// 2) Copy module contents into tmpDir/module
				modSrc := tb.Module
				if !filepath.IsAbs(modSrc) {
					modSrc = filepath.Join(filepath.Dir(p), modSrc)
				}
				modDst := filepath.Join(tmpDir, "module")
				if err := copyDir(modSrc, modDst); err != nil {
					res.Passed = false
					res.Logs = fmt.Sprintf("failed to copy module: %v", err)
					mu.Lock(); results = append(results, res); mu.Unlock()
					return
				}

				// 3) Prepare mocks
				mstore := mocks.NewMockStore()
				for _, m := range tb.Mocks { mstore.Add(m.Resource, m.Attributes) }
				if mode == "mock" {
					if err := mocks.InjectMocks(modDst, mstore); err != nil {
						res.Passed = false
						res.Logs = fmt.Sprintf("failed injecting mocks: %v", err)
						mu.Lock(); results = append(results, res); mu.Unlock()
						return
					}
				}

				// 4) Prepare test variables file
				if len(tb.Vars) > 0 {
					varsFile := filepath.Join(modDst, "terraspec_vars.tfvars.json")
					b, _ := json.Marshal(tb.Vars)
					ioutil.WriteFile(varsFile, b, 0644)
				}

				// 5) Run terraform init && plan && show -json
				out, err := terraform.ExecTerraformWithTimeout(modDst, 2*60*time.Second, "init", "-input=false")
				logs := string(out)
				if err != nil {
					res.Passed = false
					res.Logs = logs
					mu.Lock(); results = append(results, res); mu.Unlock()
					return
				}
				providersDir := filepath.Join(modDst, ".terraform", "providers")
					if _, err := os.Stat(providersDir); err == nil {
    				exec.Command("chmod", "-R", "+x", providersDir).Run()
    				exec.Command("xattr", "-dr", "com.apple.quarantine", providersDir).Run()
				}



				planOut, err := terraform.ExecTerraformWithTimeout(modDst, 2*60*time.Second, "plan", "-input=false", "-destroy=false", "-out=plan.tfplan")
				logs += "" + string(planOut)
				if err != nil {
					res.Passed = false
					res.Logs = logs
					mu.Lock(); results = append(results, res); mu.Unlock()
					return
				}

				showOut, err := terraform.ExecTerraformWithTimeout(modDst, 60*time.Second, "show", "-json", "plan.tfplan")
				logs += "" + string(showOut)
				if err != nil {
					res.Passed = false
					res.Logs = logs
					mu.Lock(); results = append(results, res); mu.Unlock()
					return
				}

				// 6) Parse outputs and planned values
				planned := map[string]interface{}{}
				if err := json.Unmarshal(showOut, &planned); err != nil {
					res.Passed = false
					res.Logs = fmt.Sprintf("failed to parse terraform json: %v", err)
					mu.Lock(); results = append(results, res); mu.Unlock()
					return
				}

				// Build a map for quick lookup for outputs: outputs.<name>
				outputsMap := map[string]interface{}{}
				if pv, ok := planned["planned_values"]; ok {
					if pvmap, ok := pv.(map[string]interface{}); ok {
						if o, ok := pvmap["outputs"]; ok {
							if om, ok := o.(map[string]interface{}); ok {
								for k, v := range om {
									if vm, ok := v.(map[string]interface{}); ok {
										outputsMap[k] = vm["value"]
									}
								}
							}
						}
					}
				}
				// Build resource map: resource.<type>.<name>.<attr>
				resourceMap := map[string]interface{}{}
				if pv, ok := planned["planned_values"]; ok {
					if pvmap, ok := pv.(map[string]interface{}); ok {
						if root, ok := pvmap["root_module"]; ok {
							collectResources(root, resourceMap)
						}
					}
				}

				// 7) Evaluate assertions. Support actual values:
				// - literal string
				// - output.<name>
				// - resource.<type>.<name>.<attr>
				ar := []assert.Result{}
				passed := true
				for _, a := range tb.Asserts {
					var actualVal interface{}
					if strings.HasPrefix(a.Actual, "output.") {
						k := strings.TrimPrefix(a.Actual, "output.")
						if v, ok := outputsMap[k]; ok {
							actualVal = v
						} else {
							actualVal = "<nil>"
						}
					} else if strings.HasPrefix(a.Actual, "resource.") {
						k := strings.TrimPrefix(a.Actual, "resource.")
						if v, ok := resourceMap[k]; ok {
							actualVal = v
						} else {
							actualVal = "<nil>"
						}
					} else {
						actualVal = a.Actual
					}

					// expected can be var.<name> or literal
					expVal := a.Expected
					if strings.HasPrefix(a.Expected, "var.") {
						k := strings.TrimPrefix(a.Expected, "var.")
						if v, ok := tb.Vars[k]; ok { expVal = v }
					}

					r := assert.Evaluate(a.Name, a.Condition, actualVal, expVal)
					ar = append(ar, r)
					if !r.Passed { passed = false }
				}

				res.Assertions = ar
				res.Passed = passed
				res.Logs = logs
				mu.Lock(); results = append(results, res); mu.Unlock()
			}(path, t)
		}
	}
	wg.Wait()

	select {
	case e := <-errorChan:
		return results, e
	default:
	}
	return results, nil
}

// copyDir copies the src directory into dst (recursive). Simplified implementation.
func copyDir(src string, dst string) error {
	info, err := os.Stat(src)
	if err != nil { return err }
	if err := os.MkdirAll(dst, info.Mode()); err != nil { return err }
	entries, err := ioutil.ReadDir(src)
	if err != nil { return err }
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil { return err }
		} else {
			if err := copyFile(srcPath, dstPath); err != nil { return err }
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil { return err }
	defer srcF.Close()
	data, err := io.ReadAll(srcF)
	if err != nil { return err }
	return ioutil.WriteFile(dst, data, 0644)
}

// collectResources traverses planned_values root_module to populate a flat map where keys are of form "<type>.<name>.<attr>" or "<type>.<name>".
func collectResources(node interface{}, out map[string]interface{}) {
	m, ok := node.(map[string]interface{})
	if !ok { return }
	// collect resources at this level
	if resources, ok := m["resources"]; ok {
		if resArr, ok := resources.([]interface{}); ok {
			for _, ri := range resArr {
				if rmap, ok := ri.(map[string]interface{}); ok {
					typeName, _ := rmap["type"].(string)
					name, _ := rmap["name"].(string)
					values, _ := rmap["values"].(map[string]interface{})
					for k, v := range values {
						key := fmt.Sprintf("%s.%s.%s", typeName, name, k)
						out[key] = v
					}
					// also store the whole values map under type.name
					out[fmt.Sprintf("%s.%s", typeName, name)] = values
				}
			}
		}
	}
	// recurse into child modules
	if child, ok := m["child_modules"]; ok {
		if arr, ok := child.([]interface{}); ok {
			for _, c := range arr { collectResources(c, out) }
			}
	}
}

