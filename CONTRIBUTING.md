# Contributing to Infractest

Thank you for your interest in contributing to Infractest! This document provides guidelines and information for contributors to help maintain code quality, consistency, and a welcoming community.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Release Process](#release-process)

## 🤝 Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow. Please be respectful, inclusive, and constructive in all interactions.

## 🚀 Getting Started

### Prerequisites

- **Go 1.25.3+**: [Install Go](https://golang.org/doc/install)
- **Terraform**: [Install Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Make** (optional): For using provided Makefile commands

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/memetics19/infractest.git
cd infractest
git remote add upstream https://github.com/memetics19/infractest.git
```

## 🛠 Development Setup

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest  # For hot reloading
```

### 2. Build the Project

```bash
# Build the binary
go build -o bin/infractest cmd/infractest/main.go

# Or use the Makefile
make build
```

### 3. Run Tests

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
make test-integration

# Run all tests
make test
```

### 4. Lint and Format

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Or use Makefile
make lint
make format
```

## 🏗 Project Structure

```
infractest/
├── cmd/
│   └── infractest/          # CLI application entry point
├── internal/
│   └── terraform/           # Internal Terraform execution logic
├── pkg/                     # Public packages
│   ├── assert/              # Assertion engine
│   ├── mocks/               # Mock injection system
│   ├── parser/              # HCL test file parser
│   ├── reporter/            # Test result reporting
│   └── runner/              # Test execution runner
├── examples/                # Example Terraform modules
├── tests/                   # Test files and examples
├── docs/                    # Documentation
├── scripts/                 # Build and utility scripts
├── .github/                 # GitHub workflows and templates
├── Makefile                 # Build automation
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── .golangci.yml            # Linter configuration
└── README.md                # Project documentation
```

## 📝 Coding Standards

### Go Code Style

We follow standard Go conventions and best practices:

#### Naming Conventions

- **Variables/Functions**: `camelCase`
- **Types/Structs**: `PascalCase`
- **Files**: `snake_case.go`
- **Constants**: `UPPER_CASE`
- **Environment Variables**: `UPPER_CASE`

#### Code Organization

```go
// Package comment
package example

// Imports (standard, third-party, local)
import (
    "fmt"
    "os"
    
    "github.com/fatih/color"
    
    "infractest/pkg/assert"
)

// Type definitions
type Example struct {
    Field string
}

// Constructor
func NewExample(field string) *Example {
    return &Example{Field: field}
}

// Methods
func (e *Example) DoSomething() error {
    // Implementation
    return nil
}
```

#### Error Handling

```go
// Always handle errors explicitly
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to execute function: %w", err)
}

// Use wrapped errors for context
if err := processData(data); err != nil {
    return fmt.Errorf("data processing failed: %w", err)
}
```

#### Logging

```go
import "log"

// Use structured logging
log.Printf("Processing test: %s", testName)
log.Printf("Test result: %+v", result)

// For errors, include context
log.Printf("Failed to execute test %s: %v", testName, err)
```

### HCL Test Files

#### File Naming

- Test files: `*.tfunittest.hcl`
- Example: `vpc_validation.tfunittest.hcl`

#### Test Structure

```hcl
# Test file header comment
# This file tests VPC module functionality

test "vpc_creation_with_valid_cidr" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
    name       = "test-vpc"
  }
  
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-12345678"
      cidr_block = "10.0.0.0/16"
    }
  }
  
  assert "vpc_cidr_matches_input" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
}
```

#### Best Practices

- Use descriptive test names
- Group related tests in the same file
- Include comments for complex test scenarios
- Use meaningful variable names
- Keep tests focused and atomic

## 🧪 Testing Guidelines

### Unit Tests

- **Location**: `*_test.go` files alongside source code
- **Coverage**: Aim for >80% code coverage
- **Naming**: `TestFunctionName` or `TestStruct_Method`

```go
func TestAssertEvaluate(t *testing.T) {
    tests := []struct {
        name      string
        condition string
        actual    interface{}
        expected  interface{}
        want      assert.Result
    }{
        {
            name:      "equals condition passes",
            condition: "equals",
            actual:    "hello",
            expected:  "hello",
            want:      assert.Result{Name: "test", Passed: true},
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := assert.Evaluate("test", tt.condition, tt.actual, tt.expected)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Evaluate() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

- **Location**: `tests/integration/`
- **Purpose**: Test end-to-end functionality
- **Execution**: Run with `make test-integration`

### Test Data

- **Location**: `testdata/`
- **Format**: Use realistic but safe test data
- **Cleanup**: Ensure tests clean up after themselves

### Mock Testing

```go
// Use interfaces for testable code
type TerraformExecutor interface {
    Execute(args []string) ([]byte, error)
}

// Mock implementation
type mockTerraformExecutor struct {
    output []byte
    err    error
}

func (m *mockTerraformExecutor) Execute(args []string) ([]byte, error) {
    return m.output, m.err
}
```

## 🔄 Pull Request Process

### Before Submitting

1. **Create a Feature Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

2. **Make Your Changes**:
   - Write clean, well-documented code
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**:
   ```bash
   make test
   make lint
   make build
   ```

4. **Commit Your Changes**:
   ```bash
   git add .
   git commit -m "feat: add new assertion condition

   - Add 'not_equals' condition to assertion engine
   - Update tests to cover new functionality
   - Update documentation with examples"
   ```

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples**:
```
feat(assert): add regex matching condition

fix(parser): handle malformed HCL test files

docs: update installation instructions

test(runner): add integration tests for live mode
```

### Pull Request Template

When creating a PR, please include:

1. **Description**: What changes were made and why
2. **Type**: Feature, Bug Fix, Documentation, etc.
3. **Testing**: How the changes were tested
4. **Breaking Changes**: Any breaking changes (if applicable)
5. **Checklist**: Ensure all items are completed

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Breaking Changes
Describe any breaking changes

## Additional Notes
Any additional information
```

### Review Process

1. **Automated Checks**: CI/CD pipeline runs automatically
2. **Code Review**: At least one maintainer review required
3. **Testing**: All tests must pass
4. **Documentation**: Update docs if needed
5. **Approval**: Maintainer approval required for merge

## 🐛 Issue Guidelines

### Before Creating an Issue

1. Search existing issues to avoid duplicates
2. Check if it's already fixed in the latest version
3. Gather relevant information

### Bug Reports

Use the bug report template and include:

- **Environment**: OS, Go version, Terraform version
- **Steps to Reproduce**: Clear, numbered steps
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Logs**: Relevant error messages or logs
- **Minimal Example**: Smallest possible test case

### Feature Requests

Use the feature request template and include:

- **Problem**: What problem does this solve?
- **Solution**: Describe your proposed solution
- **Alternatives**: Other solutions you've considered
- **Additional Context**: Any other relevant information

### Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Improvements to documentation
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention is needed
- `priority: high/medium/low`: Issue priority

## 🚀 Release Process

### Version Numbering

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. **Update Version**:
   ```bash
   # Update version in go.mod and other files
   make version VERSION=v1.2.3
   ```

2. **Update Changelog**:
   - Add new features, fixes, and breaking changes
   - Update release date

3. **Create Release**:
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

4. **Publish**:
   - GitHub automatically creates release from tag
   - Update documentation if needed

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Terraform Documentation](https://www.terraform.io/docs/)
- [HCL Language Documentation](https://github.com/hashicorp/hcl)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)

## 💡 Tips for Contributors

### Getting Help

- **GitHub Discussions**: For questions and general discussion
- **Issues**: For bug reports and feature requests
- **Code Review**: Ask questions during PR reviews

### Best Practices

1. **Start Small**: Begin with documentation or small bug fixes
2. **Ask Questions**: Don't hesitate to ask for clarification
3. **Be Patient**: Review process takes time
4. **Stay Updated**: Keep your fork up to date with upstream
5. **Test Thoroughly**: Ensure your changes work as expected

### Common Pitfalls

- **Forgetting Tests**: Always add tests for new functionality
- **Breaking Changes**: Be careful with API changes
- **Documentation**: Update docs when adding features
- **Commit Messages**: Use clear, descriptive commit messages
- **Branch Management**: Keep feature branches focused

## 🙏 Recognition

Contributors are recognized in:

- **README.md**: Major contributors
- **CHANGELOG.md**: Release notes
- **GitHub**: Contributor statistics
- **Releases**: Release notes

Thank you for contributing to Infractest! 🎉
