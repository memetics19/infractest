// ---------------------------
// pkg/assert/engine.go
// ---------------------------
package assert

import (
	"fmt"
	"strings"
)

// Simple assertion result.
type Result struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
}

// Evaluate compares actual/expected using condition.
func Evaluate(name, condition, actual, expected string) Result {
	switch strings.ToLower(condition) {
	case "equals":
		if actual == expected {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("expected %q, got %q", expected, actual)}
	case "contains":
		if strings.Contains(actual, expected) {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("%q does not contain %q", actual, expected)}
	case "matches":
		// basic substring match; regex could be added later
		if strings.Contains(actual, expected) {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("%q does not match %q", actual, expected)}
	default:
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("unknown condition %q", condition)}
	}
}
