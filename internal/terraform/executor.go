
// ---------------------------
// internal/terraform/executor.go
// ---------------------------
package terraform

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// ExecTerraform runs terraform commands with a timeout and returns combined output.
func ExecTerraform(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

// ExecTerraformWithTimeout runs terraform with a configured timeout.
func ExecTerraformWithTimeout(dir string, timeout time.Duration, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ExecTerraform(ctx, dir, args...)
}