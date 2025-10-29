# Command Line Interface Reference

Complete reference for the Infractest command line interface.

## 🚀 Basic Usage

```bash
infractest [options]
```

## 📋 Command Line Options

### `-dir` (Directory)
**Description**: Directory containing `.tfunittest.hcl` test files  
**Type**: String  
**Default**: `"tests"`  
**Example**: `infractest -dir tests`

```bash
# Run tests from specific directory
infractest -dir tests/integration

# Run tests from current directory
infractest -dir .

# Run tests from absolute path
infractest -dir /full/path/to/tests
```

### `-json` (JSON Output)
**Description**: Path to write JSON report  
**Type**: String  
**Default**: `""` (no JSON output)  
**Example**: `infractest -json report.json`

```bash
# Generate JSON report
infractest -dir tests -json test-results.json

# Generate JSON report with timestamp
infractest -dir tests -json "results-$(date +%Y%m%d-%H%M%S).json"
```

### `-mode` (Test Mode)
**Description**: Test execution mode  
**Type**: String  
**Default**: `"mock"`  
**Values**: `mock`, `live`  
**Example**: `infractest -mode live`

```bash
# Run in mock mode (default)
infractest -dir tests -mode mock

# Run in live mode
infractest -dir tests -mode live
```

## 🔧 Environment Variables

### Terraform Configuration
```bash
# Enable Terraform logging
export TF_LOG=INFO
export TF_LOG_PATH=terraform.log

# Disable interactive prompts
export TF_IN_AUTOMATION=true

# Set Terraform working directory
export TF_DATA_DIR=/tmp/terraform
```

### Cloud Provider Credentials (Live Mode)

#### AWS
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1
export AWS_SESSION_TOKEN=your_session_token  # Optional
```

#### Azure
```bash
export AZURE_CLIENT_ID=your_client_id
export AZURE_CLIENT_SECRET=your_client_secret
export AZURE_TENANT_ID=your_tenant_id
export AZURE_SUBSCRIPTION_ID=your_subscription_id
```

#### Google Cloud
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
export GOOGLE_CLOUD_PROJECT=your-project-id
```

## 📊 Exit Codes

| Code | Description |
|------|-------------|
| `0` | All tests passed |
| `1` | One or more tests failed |
| `2` | Invalid arguments or configuration error |

## 🎯 Usage Examples

### Basic Testing
```bash
# Run all tests in default directory
infractest

# Run tests from specific directory
infractest -dir tests

# Run tests and generate JSON report
infractest -dir tests -json results.json
```

### Mock Mode Testing
```bash
# Run in mock mode (default)
infractest -dir tests -mode mock

# Run in mock mode with JSON output
infractest -dir tests -mode mock -json mock-results.json
```

### Live Mode Testing
```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1

# Run in live mode
infractest -dir tests -mode live

# Run in live mode with JSON output
infractest -dir tests -mode live -json live-results.json
```

### CI/CD Integration
```bash
# GitHub Actions example
- name: Run Infrastructure Tests
  run: |
    infractest -dir tests -mode mock -json test-results.json
    infractest -dir tests -mode live -json live-results.json

# Jenkins example
pipeline {
  agent any
  stages {
    stage('Test') {
      steps {
        sh 'infractest -dir tests -mode mock -json mock-results.json'
        sh 'infractest -dir tests -mode live -json live-results.json'
      }
    }
  }
}
```

## 🔍 Debugging and Troubleshooting

### Enable Verbose Logging
```bash
# Enable Terraform debug logging
export TF_LOG=DEBUG
infractest -dir tests

# Enable Terraform trace logging
export TF_LOG=TRACE
infractest -dir tests
```

### Check Test Files
```bash
# Validate test files exist
ls -la tests/*.tfunittest.hcl

# Check test file syntax
infractest -dir tests -mode mock
```

### Verify Module Structure
```bash
# Check if modules are valid
cd modules/vpc
terraform init
terraform validate
```

## 📝 Output Formats

### Console Output
```
Running tests in mock mode...
✓ vpc_test.tfunittest.hcl::vpc_creation_with_defaults
  ✓ vpc_id_is_generated
  ✓ vpc_cidr_matches_input
  ✓ vpc_name_matches_input

Tests completed: 1 passed, 0 failed
```

### JSON Output
```json
[
  {
    "file": "vpc_test.tfunittest.hcl",
    "test": "vpc_creation_with_defaults",
    "assertions": [
      {
        "name": "vpc_id_is_generated",
        "passed": true,
        "message": ""
      },
      {
        "name": "vpc_cidr_matches_input",
        "passed": true,
        "message": ""
      }
    ],
    "passed": true,
    "logs": ""
  }
]
```

## 🚀 Performance Optimization

### Parallel Execution
Tests run in parallel by default. To control parallelism:

```bash
# Set Terraform parallelism
export TF_CLI_ARGS="-parallelism=10"
infractest -dir tests
```

### Resource Management
```bash
# Limit Terraform state locking
export TF_CLI_ARGS="-lock-timeout=30s"
infractest -dir tests

# Set Terraform plugin cache
export TF_PLUGIN_CACHE_DIR=/tmp/terraform-plugin-cache
infractest -dir tests
```

## 🔧 Advanced Configuration

### Custom Terraform Binary
```bash
# Use custom Terraform binary
export TF_BINARY_PATH=/path/to/terraform
infractest -dir tests
```

### Custom Working Directory
```bash
# Set custom working directory
export TF_DATA_DIR=/tmp/terraform-data
infractest -dir tests
```

### Plugin Cache
```bash
# Enable plugin cache
export TF_PLUGIN_CACHE_DIR=/tmp/terraform-plugin-cache
infractest -dir tests
```

## 📚 Common Use Cases

### Development Workflow
```bash
# Quick test during development
infractest -dir tests -mode mock

# Full test before commit
infractest -dir tests -mode mock -json dev-results.json
```

### CI/CD Pipeline
```bash
# Mock tests for PR validation
infractest -dir tests -mode mock -json pr-results.json

# Live tests for main branch
infractest -dir tests -mode live -json main-results.json
```

### Production Validation
```bash
# Validate production modules
infractest -dir tests/production -mode live -json prod-results.json
```

## 🆘 Error Handling

### Common Errors

#### Invalid Directory
```bash
# Error: directory not found
infractest -dir nonexistent
# Solution: Check directory path

# Error: no test files found
infractest -dir empty
# Solution: Add .tfunittest.hcl files
```

#### Terraform Errors
```bash
# Error: terraform not found
# Solution: Install Terraform

# Error: invalid module
# Solution: Validate module with terraform validate
```

#### Cloud Provider Errors
```bash
# Error: AWS credentials not found
# Solution: Set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY

# Error: insufficient permissions
# Solution: Check IAM permissions
```

## 📖 Help and Version

### Help
```bash
infractest --help
infractest -h
```

### Version
```bash
infractest --version
infractest -v
```

## 🎯 Best Practices

1. **Use descriptive test directories** - organize tests logically
2. **Generate JSON reports** - for CI/CD integration and analysis
3. **Use appropriate test modes** - mock for development, live for integration
4. **Set environment variables** - for consistent behavior
5. **Enable logging** - for debugging and troubleshooting
6. **Validate modules** - before writing tests
7. **Use parallel execution** - for faster test runs
8. **Handle errors gracefully** - check exit codes and logs

## 🔄 Integration Examples

### Makefile
```makefile
.PHONY: test test-mock test-live test-report

test: test-mock test-live

test-mock:
	infractest -dir tests -mode mock

test-live:
	infractest -dir tests -mode live

test-report:
	infractest -dir tests -json test-results.json
```

### Shell Script
```bash
#!/bin/bash
set -e

echo "Running infrastructure tests..."

# Run mock tests
echo "Running mock tests..."
infractest -dir tests -mode mock -json mock-results.json

# Run live tests if credentials are available
if [ -n "$AWS_ACCESS_KEY_ID" ]; then
  echo "Running live tests..."
  infractest -dir tests -mode live -json live-results.json
else
  echo "Skipping live tests - no AWS credentials"
fi

echo "Tests completed successfully!"
```

### Docker
```dockerfile
FROM golang:1.25.3-alpine

# Install Terraform
RUN apk add --no-cache wget unzip
RUN wget https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0_linux_amd64.zip
RUN unzip terraform_1.5.0_linux_amd64.zip
RUN mv terraform /usr/local/bin/

# Install Infractest
RUN go install github.com/memetics19/infractest/cmd/infractest@latest

WORKDIR /workspace
ENTRYPOINT ["infractest"]
```

```bash
# Run with Docker
docker run -v $(pwd):/workspace infractest -dir tests
```
