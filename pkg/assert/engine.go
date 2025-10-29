// ---------------------------
// pkg/assert/engine.go
// ---------------------------
package assert

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Simple assertion result.
type Result struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
}

// Evaluate compares actual/expected using condition. Actual/Expected can be literal
// or JSON unmarshalled values represented as interface{} converted to strings here.
func Evaluate(name, condition string, actual interface{}, expected interface{}) Result {
	actStr := fmt.Sprintf("%v", actual)
	expStr := fmt.Sprintf("%v", expected)

	switch strings.ToLower(condition) {
	case "equals":
		if actStr == expStr {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("expected %q, got %q", expStr, actStr)}
	case "contains":
		if strings.Contains(actStr, expStr) {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("%q does not contain %q", actStr, expStr)}
	case "matches":
		re, err := regexp.Compile(expStr)
		if err != nil { return Result{Name: name, Passed: false, Message: fmt.Sprintf("invalid regex %q", expStr)} }
		if re.MatchString(actStr) {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("%q does not match %q", actStr, expStr)}
	case "json_equals":
		// Attempt deep JSON compare
		var a, b interface{}
		if err := json.Unmarshal([]byte(actStr), &a); err != nil {
			return Result{Name: name, Passed: false, Message: "actual is not valid JSON"}
		}
		if err := json.Unmarshal([]byte(expStr), &b); err != nil {
			return Result{Name: name, Passed: false, Message: "expected is not valid JSON"}
		}
		if fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b) {
			return Result{Name: name, Passed: true}
		}
		return Result{Name: name, Passed: false, Message: "json structures differ"}
	default:
		return Result{Name: name, Passed: false, Message: fmt.Sprintf("unknown condition %q", condition)}
	}
}
