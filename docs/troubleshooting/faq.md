# Frequently Asked Questions

Common questions and answers about Infractest.

## 🚀 General Questions

### What is Infractest?
Infractest is a testing framework for Terraform infrastructure code that allows you to write unit tests for your Terraform modules using both mock and live execution modes.

### Why should I use Infractest?
- **Fast Testing**: Mock mode provides fast feedback during development
- **Comprehensive Testing**: Test both module logic and real cloud resources
- **CI/CD Integration**: Easy integration with GitHub Actions, Jenkins, etc.
- **Flexible Assertions**: Multiple assertion conditions for different test scenarios
- **Parallel Execution**: Tests run in parallel for faster feedback

### How is Infractest different from other Terraform testing tools?
- **Dual Modes**: Both mock and live testing in one tool
- **HCL Test Syntax**: Write tests in familiar HCL format
- **Mock Injection**: Automatic mock resource injection
- **Flexible Assertions**: Multiple assertion conditions
- **CI/CD Ready**: Built-in JSON reporting and exit codes

### What versions of Terraform are supported?
Infractest supports Terraform 0.12+ and is tested with the latest stable versions.

### What cloud providers are supported?
Infractest works with any Terraform provider, including:
- AWS
- Azure
- Google Cloud
- HashiCorp Cloud
- And many others

## 🛠 Installation Questions

### How do I install Infractest?
```bash
# Using go install (recommended)
go install github.com/memetics19/infractest/cmd/infractest@latest

# Or build from source
git clone https://github.com/memetics19/infractest.git
cd infractest
go build -o bin/infractest cmd/infractest/main.go
```

### What are the system requirements?
- Go 1.25.3 or later
- Terraform (for running tests)
- Cloud provider credentials (for live mode)

### Can I install Infractest without Go?
Yes, you can download pre-built binaries from [GitHub Releases](https://github.com/memetics19/infractest/releases).

### How do I update Infractest?
```bash
# Update to latest version
go install github.com/memetics19/infractest/cmd/infractest@latest

# Or from source
cd infractest
git pull origin main
go build -o bin/infractest cmd/infractest/main.go
```

## 🧪 Testing Questions

### What is the difference between mock and live mode?
- **Mock Mode**: Fast unit testing with injected mock resources
- **Live Mode**: Integration testing with real cloud resources

### When should I use mock mode vs live mode?
- **Mock Mode**: Development, unit testing, fast feedback
- **Live Mode**: Integration testing, end-to-end validation, CI/CD

### How do I write my first test?
```hcl
test "my_first_test" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
  }
  
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-12345678"
      cidr_block = "10.0.0.0/16"
    }
  }
  
  assert "vpc_id_is_generated" {
    actual    = "output.vpc_id"
    expected  = "vpc-12345678"
    condition = "equals"
  }
}
```

### How do I test modules with complex dependencies?
Use mock resources to simulate dependencies:

```hcl
test "complex_module_test" {
  module = "../modules/complex"
  
  # Mock all dependencies
  mock "aws_vpc.main" {
    attributes = {
      id = "vpc-12345678"
    }
  }
  
  mock "aws_subnet.public" {
    attributes = {
      id     = "subnet-12345678"
      vpc_id = "vpc-12345678"
    }
  }
  
  # Test the module
  assert "module_works" {
    actual    = "output.result"
    expected  = "success"
    condition = "equals"
  }
}
```

### How do I test error conditions?
```hcl
test "error_condition_test" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "invalid-cidr"  # Invalid input
  }
  
  # Test that error is handled
  assert "error_is_handled" {
    actual    = "output.error_message"
    expected  = "invalid CIDR"
    condition = "contains"
  }
}
```

## 🎭 Mock Questions

### How do I mock complex resources?
```hcl
mock "aws_instance.web" {
  attributes = {
    id               = "i-12345678"
    ami              = "ami-12345678"
    instance_type    = "t3.micro"
    subnet_id        = "subnet-12345678"
    security_groups  = ["sg-12345678"]
    private_ip       = "10.0.1.100"
    public_ip        = "203.0.113.100"
    tags = {
      Name        = "web-server"
      Environment = "test"
    }
  }
}
```

### Do I need to mock all resources?
No, only mock the resources that your module directly references. Infractest will handle the rest.

### How do I mock resources with complex attributes?
Use nested maps and lists:

```hcl
mock "aws_security_group.web" {
  attributes = {
    id          = "sg-12345678"
    name        = "web-sg"
    description = "Security group for web servers"
    vpc_id      = "vpc-12345678"
    ingress = [
      {
        from_port   = 80
        to_port     = 80
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
      }
    ]
    egress = [
      {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
      }
    ]
    tags = {
      Name = "web-security-group"
    }
  }
}
```

## ✅ Assertion Questions

### What assertion conditions are available?
- **`equals`**: Exact string match
- **`contains`**: Substring match
- **`matches`**: Regex pattern match
- **`json_equals`**: Deep JSON comparison

### How do I test complex outputs?
Use `json_equals` for complex structures:

```hcl
assert "complex_output_matches" {
  actual    = "output.complex_data"
  expected  = '{"key1":"value1","key2":["item1","item2"]}'
  condition = "json_equals"
}
```

### How do I test list outputs?
Use `contains` for list elements:

```hcl
assert "list_contains_item" {
  actual    = "output.subnet_ids"
  expected  = "subnet-"
  condition = "contains"
}
```

### How do I test regex patterns?
Use `matches` condition:

```hcl
assert "valid_vpc_id" {
  actual    = "output.vpc_id"
  expected  = "vpc-[0-9a-f]{8}"
  condition = "matches"
}
```

## 🌐 Live Mode Questions

### How do I set up cloud provider credentials?
```bash
# AWS
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1

# Azure
export AZURE_CLIENT_ID=your_client_id
export AZURE_CLIENT_SECRET=your_secret
export AZURE_TENANT_ID=your_tenant_id
export AZURE_SUBSCRIPTION_ID=your_subscription_id

# Google Cloud
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
export GOOGLE_CLOUD_PROJECT=your-project-id
```

### How do I avoid resource conflicts in live mode?
- Use unique resource names
- Add random suffixes or timestamps
- Clean up resources after testing
- Use separate test environments

### How do I test with different regions?
```bash
# Set different region
export AWS_DEFAULT_REGION=eu-west-1
infractest -dir tests -mode live
```

### How do I handle costs in live mode?
- Use small instance types
- Test in development environments
- Clean up resources after testing
- Use spot instances where possible

## 🔧 Configuration Questions

### How do I configure Terraform logging?
```bash
# Enable debug logging
export TF_LOG=DEBUG
infractest -dir tests

# Set log file
export TF_LOG_PATH=terraform.log
infractest -dir tests
```

### How do I set custom Terraform options?
```bash
# Set parallelism
export TF_CLI_ARGS="-parallelism=10"
infractest -dir tests

# Set lock timeout
export TF_CLI_ARGS="-lock-timeout=30s"
infractest -dir tests
```

### How do I use custom Terraform binary?
```bash
# Set custom binary path
export TF_BINARY_PATH=/path/to/terraform
infractest -dir tests
```

## 🚀 Performance Questions

### How can I make tests run faster?
- Use mock mode for development
- Use live mode only for integration testing
- Optimize Terraform parallelism
- Use smaller test datasets

### How do I run tests in parallel?
Tests run in parallel by default. To control parallelism:

```bash
# Set Terraform parallelism
export TF_CLI_ARGS="-parallelism=10"
infractest -dir tests
```

### How do I handle memory issues?
```bash
# Increase memory limits
export GOGC=100
infractest -dir tests

# Reduce parallelism
export TF_CLI_ARGS="-parallelism=5"
infractest -dir tests
```

## 🔍 Debugging Questions

### How do I debug test failures?
```bash
# Enable debug logging
export TF_LOG=DEBUG
infractest -dir tests

# Check JSON output
infractest -dir tests -json results.json
cat results.json | jq '.'
```

### How do I validate my module before testing?
```bash
# Validate module
cd modules/your-module
terraform init
terraform validate
terraform plan
```

### How do I check test file syntax?
```bash
# Run tests to check syntax
infractest -dir tests -mode mock

# Use HCL linter if available
```

## 🔄 CI/CD Questions

### How do I integrate with GitHub Actions?
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
    - name: Install Infractest
      run: go install github.com/memetics19/infractest/cmd/infractest@latest
    - name: Run Tests
      run: infractest -dir tests -json results.json
    - name: Upload Results
      uses: actions/upload-artifact@v3
      with:
        name: test-results
        path: results.json
```

### How do I integrate with Jenkins?
```groovy
pipeline {
  agent any
  stages {
    stage('Test') {
      steps {
        sh 'go install github.com/memetics19/infractest/cmd/infractest@latest'
        sh 'infractest -dir tests -json results.json'
      }
    }
  }
  post {
    always {
      archiveArtifacts artifacts: 'results.json'
    }
  }
}
```

### How do I generate test reports?
```bash
# Generate JSON report
infractest -dir tests -json results.json

# Process JSON report
cat results.json | jq '.'
```

## 🆘 Support Questions

### Where can I get help?
- **GitHub Issues**: [Report bugs and request features](https://github.com/memetics19/infractest/issues)
- **GitHub Discussions**: [Ask questions and discuss](https://github.com/memetics19/infractest/discussions)
- **Documentation**: [Comprehensive guides](https://github.com/memetics19/infractest/docs)

### How do I report a bug?
1. Check existing issues
2. Create new issue with:
   - Error messages
   - Steps to reproduce
   - Environment details
   - Logs and output

### How do I request a feature?
1. Check existing issues and discussions
2. Create new issue with:
   - Feature description
   - Use case
   - Proposed solution
   - Benefits

### How do I contribute?
See [CONTRIBUTING.md](https://github.com/memetics19/infractest/CONTRIBUTING.md) for guidelines on:
- Development setup
- Coding standards
- Testing guidelines
- Pull request process

## 📚 Additional Resources

- **[Installation Guide](../installation/README.md)** - Detailed installation instructions
- **[Quick Start Guide](../usage/quick-start.md)** - Get started quickly
- **[Writing Tests](../usage/writing-tests.md)** - Learn to write effective tests
- **[Basic Examples](../examples/basic.md)** - Practical examples
- **[API Reference](../api/cli.md)** - Complete CLI reference
- **[Troubleshooting](common-issues.md)** - Solve common problems
