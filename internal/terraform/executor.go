// ---------------------------
// internal/terraform/executor.go
// ---------------------------
package terraform

import (
	"context"
	"os/exec"
	"time"
)

// ExecTerraform runs terraform commands with a timeout. This is a small helper
// to later support invoking `terraform plan` and injecting mocks if needed.
func ExecTerraform(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}