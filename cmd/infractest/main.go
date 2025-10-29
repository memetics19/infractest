
// ---------------------------
// cmd/terraspec/main.go
// ---------------------------
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/memetics19/infractest/pkg/runner"
	"github.com/memetics19/infractest/pkg/reporter"
)

func main() {
	var dir string
	var jsonOut string
	flag.StringVar(&dir, "dir", "tests", "directory containing .tfunittest.hcl test files")
	flag.StringVar(&jsonOut, "json", "", "path to write JSON report (optional)")
	flag.Parse()

	abs, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println("invalid test directory:", err)
		os.Exit(2)
	}

	results, err := runner.RunDirectory(abs)
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
    actual    = "10.0.0.0/16"
    expected  = "10.0.0.0/16"
    condition = "equals"
  }
}
*/

// ---------------------------
// Makefile
// ---------------------------
/*****
BIN := bin/terraspec

build:
	@echo "Building TerraSpec..."
	@mkdir -p bin
	@go build -o $(BIN) ./cmd/terraspec

run:
	@$(BIN) -dir=tests

fmt:
	@go fmt ./...

lint:
	@golangci-lint run

clean:
	@rm -rf bin
*****/

// ---------------------------
// README.md (brief)
// ---------------------------
/*
TerraSpec — Fast, Mockable Unit Testing for Terraform Modules

Quickstart (macOS):
  brew install go terraform make
  git clone https://github.com/<you>/terraspec.git
  cd terraspec
  go mod tidy
  make build
  ./bin/terraspec -dir=tests

Project layout: see comments at top of this file.
*/
