# Infractest

[![Go Version](https://img.shields.io/badge/go-1.25.3-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

**Infractest** is a powerful testing framework for Terraform infrastructure code that enables you to write comprehensive unit tests for your Terraform modules using both mock and live execution modes. Built with Go, it provides fast, reliable testing capabilities for infrastructure-as-code projects.

## 🚀 Features

- **Dual Testing Modes**: Test your Terraform modules with both mock and live execution
- **Mock Injection**: Automatically inject mock resources to test module logic without real cloud resources
- **Flexible Assertions**: Support for multiple assertion conditions (equals, contains, matches, json_equals)
- **Parallel Execution**: Run multiple tests concurrently for faster feedback
- **JSON Reporting**: Generate detailed JSON reports for CI/CD integration
- **HCL Test Syntax**: Write tests using familiar HCL syntax
- **Sandboxed Execution**: Each test runs in an isolated environment
- **Variable Injection**: Pass test variables to your modules seamlessly

## 📋 Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Test Syntax](#test-syntax)
- [Assertion Conditions](#assertion-conditions)
- [Testing Modes](#testing-modes)
- [Examples](#examples)
- [Configuration](#configuration)
- [CI/CD Integration](#cicd-integration)
- [Contributing](#contributing)
- [License](#license)

## 🛠 Installation

### Prerequisites

- Go 1.25.3 or later
- Terraform (for live mode testing)

### Install from Source

```bash
# Clone the repository
git clone https://github.com/memetics19/infractest.git
cd infractest

# Build the binary
go build -o bin/infractest cmd/infractest/main.go

# Make it executable
chmod +x bin/infractest

# Add to PATH (optional)
export PATH=$PATH:$(pwd)/bin
```

### Install via Go Install

```bash
go install github.com/memetics19/infractest/cmd/infractest@latest
```

## 🚀 Quick Start

1. **Create a test file** (`tests/example.tfunittest.hcl`):

```hcl
test "vpc cidr validation" {
  module = "../examples/vpc"

  vars = {
    cidr_block = "10.0.0.0/16"
    mock_mode  = true
  }

  assert "cidr matches variable" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
}
```

2. **Run the tests**:

```bash
# Run in mock mode (default)
infractest -dir tests

# Run in live mode
infractest -dir tests -mode live

# Generate JSON report
infractest -dir tests -json report.json
```

## 📝 Test Syntax

### Basic Test Structure

```hcl
test "test_name" {
  module = "path/to/terraform/module"
  
  vars = {
    variable_name = "value"
  }
  
  mock "resource_type.name" {
    attributes = {
      attribute_name = "value"
    }
  }
  
  assert "assertion_name" {
    actual    = "output.output_name"
    expected  = "var.variable_name"
    condition = "equals"
  }
}
```

### Test Block Components

- **`module`**: Path to the Terraform module to test (relative or absolute)
- **`vars`**: Variables to pass to the module during testing
- **`mock`**: Mock resources to inject (mock mode only)
- **`assert`**: Assertions to validate against module outputs

## 🔍 Assertion Conditions

| Condition | Description | Example |
|-----------|-------------|---------|
| `equals` | Exact string match | `"hello"` equals `"hello"` |
| `contains` | Substring match | `"hello world"` contains `"world"` |
| `matches` | Regex pattern match | `"test123"` matches `"test\\d+"` |
| `json_equals` | Deep JSON comparison | `{"a": 1}` json_equals `{"a": 1}` |

### Reference Types

- **`output.name`**: Reference module outputs
- **`var.name`**: Reference module variables
- **`resource.type.name.attribute`**: Reference resource attributes (live mode)
- **Literal values**: Direct string/number values

## 🎯 Testing Modes

### Mock Mode (Default)

- **Purpose**: Fast unit testing without cloud resources
- **Execution**: Injects mock resources and runs `terraform plan`
- **Use Case**: Testing module logic, variable handling, and output generation
- **Performance**: Very fast, no external dependencies

```bash
infractest -dir tests -mode mock
```

### Live Mode

- **Purpose**: Integration testing with real cloud resources
- **Execution**: Runs `terraform init`, `plan`, and `apply` with real providers
- **Use Case**: End-to-end testing, provider compatibility, real resource validation
- **Performance**: Slower, requires cloud credentials

```bash
infractest -dir tests -mode live
```

## 📚 Examples

### VPC Module Testing

```hcl
test "vpc creation with custom cidr" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "172.16.0.0/16"
    name       = "test-vpc"
  }
  
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-12345678"
      cidr_block = "172.16.0.0/16"
      tags = {
        Name = "test-vpc"
      }
    }
  }
  
  assert "vpc cidr is correct" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
  
  assert "vpc has correct name tag" {
    actual    = "output.vpc_tags"
    expected  = "{\"Name\":\"test-vpc\"}"
    condition = "json_equals"
  }
}
```

### Security Group Testing

```hcl
test "security group allows ssh" {
  module = "../modules/security-group"
  
  vars = {
    allow_ssh = true
    vpc_id    = "vpc-12345678"
  }
  
  assert "ssh port is open" {
    actual    = "output.security_group_rules"
    expected  = "22"
    condition = "contains"
  }
}
```

## ⚙️ Configuration

### Command Line Options

```bash
infractest [options]

Options:
  -dir string
        directory containing .tfunittest.hcl test files (default "tests")
  -json string
        path to write JSON report (optional)
  -mode string
        test mode: mock | live (default "mock")
```

### Environment Variables

```bash
# Terraform configuration
export TF_LOG=INFO
export TF_LOG_PATH=terraform.log

# AWS credentials (for live mode)
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1
```

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Infrastructure Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.25.3'
    
    - name: Setup Terraform
      uses: hashicorp/setup-terraform@v2
      with:
        terraform_version: 1.5.0
    
    - name: Install Infractest
      run: go install github.com/memetics19/infractest/cmd/infractest@latest
    
    - name: Run Mock Tests
      run: infractest -dir tests -mode mock
    
    - name: Run Live Tests
      if: github.event_name == 'pull_request'
      env:
        AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
        AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
      run: infractest -dir tests -mode live -json test-results.json
    
    - name: Upload Test Results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: test-results
        path: test-results.json
```

### Jenkins Pipeline Example

```groovy
pipeline {
    agent any
    
    stages {
        stage('Setup') {
            steps {
                sh 'go version'
                sh 'terraform version'
                sh 'go install github.com/memetics19/infractest/cmd/infractest@latest'
            }
        }
        
        stage('Mock Tests') {
            steps {
                sh 'infractest -dir tests -mode mock'
            }
        }
        
        stage('Live Tests') {
            when {
                branch 'main'
            }
            steps {
                withCredentials([[
                    $class: 'AmazonWebServicesCredentialsBinding',
                    credentialsId: 'aws-credentials'
                ]]) {
                    sh 'infractest -dir tests -mode live -json live-test-results.json'
                }
            }
        }
    }
    
    post {
        always {
            archiveArtifacts artifacts: '*.json', fingerprint: true
            publishTestResults testResultsPattern: '*.json'
        }
    }
}
```

## 🏗 Project Structure

```
infractest/
├── cmd/
│   └── infractest/          # Main CLI application
├── internal/
│   └── terraform/           # Terraform execution logic
├── pkg/
│   ├── assert/              # Assertion engine
│   ├── mocks/               # Mock injection system
│   ├── parser/              # HCL test file parser
│   ├── reporter/            # Test result reporting
│   └── runner/              # Test execution runner
├── examples/                # Example Terraform modules
├── tests/                   # Test files
└── bin/                     # Built binaries
```

## 🤝 Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for details on:

- Development setup
- Coding standards
- Testing guidelines
- Pull request process

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [HashiCorp](https://www.hashicorp.com/) for Terraform
- [HCL](https://github.com/hashicorp/hcl) for the configuration language
- The Go community for excellent tooling and libraries

## 📞 Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/memetics19/infractest/issues)
- 💬 [Discussions](https://github.com/memetics19/infractest/discussions)

---

**Made with ❤️ for the Infrastructure as Code community**
