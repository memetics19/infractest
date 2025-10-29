// ---------------------------
// pkg/reporter/reporter.go
// ---------------------------
package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/memetics19/infractest/pkg/runner"
	"github.com/fatih/color"
)

func StdoutReporter(results []runner.TestCaseResult) {
	for _, r := range results {
		if r.Passed {
			color.Green("✔ %s :: %s", r.File, r.Test)
		} else {
			color.Red("✘ %s :: %s", r.File, r.Test)
		}
		for _, a := range r.Assertions {
			if a.Passed {
				fmt.Printf("    - %s: PASS
", a.Name)
			} else {
				fmt.Printf("    - %s: FAIL — %s
", a.Name, a.Message)
			}
		}
	}
}

func WriteJSON(results []runner.TestCaseResult, path string) error {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil { return err }
	return os.WriteFile(path, b, 0644)
}