// ---------------------------
// pkg/parser/parser.go
// ---------------------------
package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// TestFile represents a collection of test blocks parsed from an HCL file.
type TestFile struct {
	Tests []TestBlock `hcl:"test,block"`
}

// TestBlock describes a single test case.
type TestBlock struct {
	Name    string            `hcl:",label"`
	Module  string            `hcl:"module,attr"`
	Vars    map[string]string `hcl:"vars,attr"`
	Mocks   []MockBlock       `hcl:"mock,block"`
	Asserts []AssertBlock     `hcl:"assert,block"`
}

// MockBlock defines a resource mock.
type MockBlock struct {
	Resource   string            `hcl:",label"`
	Attributes map[string]string `hcl:"attributes,attr"`
}

// AssertBlock defines an assertion.
type AssertBlock struct {
	Name      string `hcl:",label"`
	Actual    string `hcl:"actual,attr"`
	Expected  string `hcl:"expected,attr"`
	Condition string `hcl:"condition,attr"`
}

// ParseFile parses a single .tfunittest.hcl file.
func ParseFile(path string) (*TestFile, error) {
	var tf TestFile
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	if err := hclsimple.DecodeFile(path, nil, &tf); err != nil {
		return nil, err
	}
	return &tf, nil
}

// ParseDirectory finds and parses all .tfunittest.hcl files in dir.
func ParseDirectory(dir string) (map[string]*TestFile, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.tfunittest.hcl"))
	if err != nil {
		return nil, err
	}
	out := make(map[string]*TestFile)
	for _, f := range files {
		pf, err := ParseFile(f)
		if err != nil {
			return nil, err
		}
		out[f] = pf
	}
	return out, nil
}
