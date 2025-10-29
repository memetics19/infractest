// TerraSpec — Upgraded Runner with Terraform execution & mock injection
// ======================================================
// This file contains a set of Go source files concatenated for convenience.
// It implements sandboxed module copying, mock injection (as generated .tf files),
// terraform init/plan/show (-json) invocation, and JSON parsing for outputs
// to evaluate assertions that reference `output.<name>` or literal values.

// ---------------------------
// cmd/terraspec/main.go
// ---------------------------
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"infractest/pkg/runner"
	"infractest/pkg/reporter"
)

func main() {
	var dir string
	var jsonOut string
	var mode string
	flag.StringVar(&dir, "dir", "tests", "directory containing .tfunittest.hcl test files")
	flag.StringVar(&jsonOut, "json", "", "path to write JSON report (optional)")
	flag.StringVar(&mode, "mode", "mock", "test mode: mock | live")
	flag.Parse()

	abs, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println("invalid test directory:", err)
		os.Exit(2)
	}

	results, err := runner.RunDirectory(abs, mode)
	if err != nil {
		reporter.StdoutReporter(results)
		fmt.Println("run errors:", err)
		os.Exit(1)
	}

	reporter.StdoutReporter(results)
	if jsonOut != "" {
		if err := reporter.WriteJSON(results, jsonOut); err != nil {
			fmt.Println("failed to write json:", err)
			os.Exit(1)
		}
	}
}





// ---------------------------
// tests/example.tfunittest.hcl
// ---------------------------
/*
Place in tests/example.tfunittest.hcl

test "vpc cidr validation" {
  module = "../examples/vpc"

  vars = {
    cidr_block = "10.0.0.0/16"
  }

  mock "aws_vpc.main" {
    attributes = {
      id = "vpc-123"
      cidr_block = "10.0.0.0/16"
    }
  }

  assert "cidr matches variable" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
}
*/
