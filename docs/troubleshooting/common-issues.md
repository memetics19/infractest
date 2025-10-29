# Common Issues and Solutions

This guide helps you troubleshoot common issues when using Infractest.

## 🚨 Installation Issues

### Go Not Found
**Error**: `go: command not found`

**Solution**:
```bash
# Install Go from https://golang.org/dl/
# Or use package manager
brew install go                    # macOS
sudo apt install golang-go        # Ubuntu
sudo yum install golang           # CentOS/RHEL
```

### Permission Denied
**Error**: `permission denied` when installing

**Solution**:
```bash
# Install to user directory
go install github.com/memetics19/infractest/cmd/infractest@latest

# Or use sudo (not recommended)
sudo go install github.com/memetics19/infractest/cmd/infractest@latest
```

### Module Not Found
**Error**: `module not found` when installing

**Solution**:
```bash
# Ensure Go modules are enabled
export GO111MODULE=on
go install github.com/memetics19/infractest/cmd/infractest@latest

# Or use go get
go get github.com/memetics19/infractest/cmd/infractest@latest
```

## 🔧 Terraform Issues

### Terraform Not Found
**Error**: `terraform: command not found`

**Solution**:
```bash
# Install Terraform
# See: https://learn.hashicorp.com/tutorials/terraform/install-cli

# Verify installation
terraform version
```

### Invalid Module
**Error**: `Error: Invalid module`

**Solution**:
```bash
# Validate the module
cd modules/your-module
terraform init
terraform validate

# Check for syntax errors
terraform fmt -check
```

### Provider Issues
**Error**: `Error: Failed to query available provider packages`

**Solution**:
```bash
# Initialize Terraform
terraform init

# Update provider versions
terraform init -upgrade

# Clear provider cache
rm -rf .terraform
terraform init
```

## 🎭 Mock Mode Issues

### Missing Mock Resources
**Error**: `Error: Resource not found` or test fails

**Solution**:
```hcl
# Ensure all required resources are mocked
test "my_test" {
  module = "../modules/vpc"
  
  # Mock all resources used by the module
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-12345678"
      cidr_block = "10.0.0.0/16"
    }
  }
  
  # Add more mocks as needed
  mock "aws_subnet.public" {
    attributes = {
      id     = "subnet-12345678"
      vpc_id = "vpc-12345678"
    }
  }
}
```

### Incorrect Mock Attributes
**Error**: `Error: Invalid attribute` or assertion fails

**Solution**:
```hcl
# Use correct attribute names and types
mock "aws_vpc.main" {
  attributes = {
    id                = "vpc-12345678"        # String
    cidr_block        = "10.0.0.0/16"         # String
    enable_dns_hostnames = true               # Boolean
    enable_dns_support   = true               # Boolean
    tags = {                                  # Map
      Name = "test-vpc"
    }
  }
}
```

### Mock Resource Format
**Error**: `Error: Invalid mock resource format`

**Solution**:
```hcl
# Use correct format: resource_type.name
mock "aws_vpc.main" {          # Correct
  attributes = { ... }
}

# Avoid these formats:
mock "aws_vpc" {               # Missing name
mock "vpc.main" {              # Missing type
mock "aws_vpc.main.attributes" { # Too specific
```

## 🌐 Live Mode Issues

### Missing Credentials
**Error**: `Error: No valid credential sources found`

**Solution**:
```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1

# Or use AWS CLI
aws configure

# Verify credentials
aws sts get-caller-identity
```

### Insufficient Permissions
**Error**: `Error: Access Denied` or `Error: Forbidden`

**Solution**:
```bash
# Check IAM permissions
aws iam get-user
aws iam list-attached-user-policies --user-name your-username

# Ensure user has required permissions:
# - ec2:CreateVpc
# - ec2:CreateSubnet
# - ec2:CreateSecurityGroup
# - etc.
```

### Resource Conflicts
**Error**: `Error: Resource already exists`

**Solution**:
```bash
# Use unique resource names
# Add random suffixes or timestamps
# Clean up existing resources
aws ec2 describe-vpcs --filters "Name=tag:Name,Values=test-vpc"
aws ec2 delete-vpc --vpc-id vpc-12345678
```

### Region Issues
**Error**: `Error: Invalid region`

**Solution**:
```bash
# Set correct region
export AWS_DEFAULT_REGION=us-east-1

# Verify region is valid
aws ec2 describe-regions --region-names us-east-1
```

## 🧪 Test File Issues

### Invalid Test File
**Error**: `Error: Invalid test file` or parsing errors

**Solution**:
```hcl
# Check HCL syntax
# Ensure proper indentation
# Use correct block structure

test "my_test" {                    # Correct
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
  }
  
  assert "test_name" {
    actual    = "output.vpc_id"
    expected  = "vpc-12345678"
    condition = "equals"
  }
}
```

### Missing Test Files
**Error**: `Error: No test files found`

**Solution**:
```bash
# Check directory structure
ls -la tests/
ls -la tests/*.tfunittest.hcl

# Ensure test files have correct extension
mv tests/my_test.tf tests/my_test.tfunittest.hcl
```

### Module Path Issues
**Error**: `Error: Module not found`

**Solution**:
```hcl
# Use correct relative paths
test "my_test" {
  module = "../modules/vpc"        # Correct
  # module = "modules/vpc"         # Wrong if tests/ is parent
  # module = "../../modules/vpc"   # Wrong if tests/ is grandparent
}
```

## ✅ Assertion Issues

### Assertion Failures
**Error**: `Assertion failed: expected X, got Y`

**Solution**:
```hcl
# Check actual vs expected values
assert "my_assertion" {
  actual    = "output.vpc_id"           # Check this value
  expected  = "vpc-12345678"            # Check this value
  condition = "equals"
}

# Use correct reference format
# output.name for outputs
# var.name for variables
# resource.type.name.attribute for resources (live mode)
```

### Wrong Reference Format
**Error**: `Error: Invalid reference`

**Solution**:
```hcl
# Use correct reference format
actual = "output.vpc_id"              # Correct
actual = "var.cidr_block"             # Correct
actual = "resource.aws_vpc.main.id"   # Correct (live mode)

# Avoid these formats:
actual = "outputs.vpc_id"             # Wrong
actual = "variables.cidr_block"       # Wrong
actual = "aws_vpc.main.id"            # Wrong
```

### Condition Issues
**Error**: `Error: Invalid condition`

**Solution**:
```hcl
# Use valid conditions
condition = "equals"        # String equality
condition = "contains"      # Substring match
condition = "matches"       # Regex match
condition = "json_equals"   # JSON comparison

# Avoid invalid conditions:
condition = "equal"         # Wrong
condition = "match"         # Wrong
condition = "json_equal"    # Wrong
```

## 🔍 Debugging Issues

### Enable Debug Logging
**Problem**: Not enough information to debug

**Solution**:
```bash
# Enable Terraform debug logging
export TF_LOG=DEBUG
infractest -dir tests

# Enable Terraform trace logging
export TF_LOG=TRACE
infractest -dir tests

# Check JSON output
infractest -dir tests -json results.json
cat results.json | jq '.'
```

### Check Module Structure
**Problem**: Module validation issues

**Solution**:
```bash
# Validate module
cd modules/your-module
terraform init
terraform validate
terraform plan

# Check for syntax errors
terraform fmt -check
terraform fmt -diff
```

### Verify Test Files
**Problem**: Test file issues

**Solution**:
```bash
# Check test file syntax
infractest -dir tests -mode mock

# Validate HCL syntax
# Use a HCL linter or validator
```

## 🚀 Performance Issues

### Slow Test Execution
**Problem**: Tests take too long to run

**Solution**:
```bash
# Use mock mode for development
infractest -dir tests -mode mock

# Use live mode only for integration testing
infractest -dir tests -mode live

# Optimize Terraform parallelism
export TF_CLI_ARGS="-parallelism=10"
infractest -dir tests
```

### Memory Issues
**Problem**: Out of memory errors

**Solution**:
```bash
# Increase memory limits
export GOGC=100
infractest -dir tests

# Use smaller test datasets
# Reduce parallel execution
export TF_CLI_ARGS="-parallelism=5"
infractest -dir tests
```

### Network Issues
**Problem**: Network timeouts or connectivity issues

**Solution**:
```bash
# Check network connectivity
ping google.com
nslookup registry.terraform.io

# Use proxy if needed
export HTTP_PROXY=http://proxy:port
export HTTPS_PROXY=http://proxy:port
infractest -dir tests
```

## 🔧 Configuration Issues

### Environment Variables
**Problem**: Environment variables not working

**Solution**:
```bash
# Check environment variables
env | grep TF_
env | grep AWS_

# Set variables explicitly
export TF_LOG=INFO
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
infractest -dir tests
```

### Working Directory
**Problem**: Wrong working directory

**Solution**:
```bash
# Run from correct directory
cd /path/to/your/project
infractest -dir tests

# Use absolute paths
infractest -dir /full/path/to/tests
```

### File Permissions
**Problem**: Permission denied errors

**Solution**:
```bash
# Check file permissions
ls -la tests/
ls -la modules/

# Fix permissions if needed
chmod 755 tests/
chmod 755 modules/
chmod 644 tests/*.tfunittest.hcl
```

## 🆘 Getting Help

### Check Logs
```bash
# Enable verbose logging
export TF_LOG=DEBUG
infractest -dir tests 2>&1 | tee debug.log

# Check JSON output
infractest -dir tests -json results.json
cat results.json | jq '.'
```

### Verify Installation
```bash
# Check Infractest version
infractest --version

# Check Go version
go version

# Check Terraform version
terraform version
```

### Test with Minimal Example
```hcl
# Create minimal test
test "minimal_test" {
  module = "."
  
  assert "test_passes" {
    actual    = "output.test"
    expected  = "test"
    condition = "equals"
  }
}
```

### Report Issues
If you can't resolve the issue:

1. **Check existing issues**: [GitHub Issues](https://github.com/memetics19/infractest/issues)
2. **Create new issue**: Include error messages, logs, and steps to reproduce
3. **Ask questions**: [GitHub Discussions](https://github.com/memetics19/infractest/discussions)

## 📚 Additional Resources

- **[Installation Guide](../installation/README.md)** - Installation troubleshooting
- **[Writing Tests](writing-tests.md)** - Test writing best practices
- **[Test Modes](test-modes.md)** - Understanding test modes
- **[Assertions](assertions.md)** - Assertion troubleshooting
- **[FAQ](faq.md)** - Frequently asked questions
