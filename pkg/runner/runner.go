// ---------------------------
// pkg/runner/runner.go
// ---------------------------
package runner

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/memetics19/infractest/pkg/assert"
	"github.com/memetics19/infractest/pkg/mocks"
	"github.com/memetics19/infractest/pkg/parser"
)

// TestCaseResult captures a single test's results.
type TestCaseResult struct {
	File       string                    `json:"file"`
	Test       string                    `json:"test"`
	Assertions []assert.Result          `json:"assertions"`
	Passed     bool                      `json:"passed"`
}

// RunDirectory parses and runs tests in a directory.
func RunDirectory(dir string) ([]TestCaseResult, error) {
	parsed, err := parser.ParseDirectory(dir)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []TestCaseResult{}

	for path, tf := range parsed {
		for _, t := range tf.Tests {
			wg.Add(1)
			go func(p string, tb parser.TestBlock) {
				defer wg.Done()
				mstore := mocks.NewMockStore()
				for _, m := range tb.Mocks {
					mstore.Add(m.Resource, m.Attributes)
				}

				// Evaluate asserts — in this MVP we treat actual/expected as literal strings.
				ar := []assert.Result{}
				passed := true
				for _, a := range tb.Asserts {
					r := assert.Evaluate(a.Name, a.Condition, a.Actual, a.Expected)
					ar = append(ar, r)
					if !r.Passed { passed = false }
				}

				res := TestCaseResult{File: filepath.Base(p), Test: tb.Name, Assertions: ar, Passed: passed}
				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}(path, t)
		}
	}
	wg.Wait()
	return results, nil
}