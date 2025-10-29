# Basic Concepts

Understanding the core concepts of Infractest will help you write better tests and use the tool more effectively.

## 🎯 What is Infractest?

Infractest is a testing framework specifically designed for Terraform infrastructure code. It allows you to:

- **Test Terraform modules** without deploying real cloud resources
- **Validate module logic** using mock resources
- **Ensure consistency** between expected and actual outputs
- **Integrate testing** into your CI/CD pipeline

## 🏗 Core Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Test Files    │───▶│   Infractest     │───▶│  Test Results   │
│ (.tfunittest.hcl)│    │     Runner       │    │   (JSON/CLI)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │  Terraform       │
                       │  Execution       │
                       │  (Mock/Live)     │
                       └──────────────────┘
```

## 🧪 Test Components

### 1. Test Block
The main container for a test case:

```hcl
test "descriptive_test_name" {
  # Test configuration
}
```

**Key Points**:
- Each test is independent
- Tests run in parallel by default
- Use descriptive names for clarity

### 2. Module Reference
Points to the Terraform module to test:

```hcl
test "my_test" {
  module = "../modules/vpc"  # Relative or absolute path
}
```

**Supported Paths**:
- Relative: `../modules/vpc`
- Absolute: `/full/path/to/module`
- Current directory: `"."`

### 3. Variables
Pass input variables to the module:

```hcl
test "my_test" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
    name       = "test-vpc"
    tags = {
      Environment = "test"
      Owner       = "team"
    }
  }
}
```

**Variable Types**:
- **Strings**: `"value"`
- **Numbers**: `42`, `3.14`
- **Booleans**: `true`, `false`
- **Lists**: `["item1", "item2"]`
- **Maps**: `{ key = "value" }`

### 4. Mock Resources
Define mock versions of cloud resources:

```hcl
test "my_test" {
  module = "../modules/vpc"
  
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-12345678"
      cidr_block = "10.0.0.0/16"
      tags = {
        Name = "test-vpc"
      }
    }
  }
}
```

**Mock Resource Format**:
- **Resource Type**: `aws_vpc.main`
- **Attributes**: Key-value pairs of resource attributes
- **Nested Values**: Support for complex nested structures

### 5. Assertions
Validate expected behavior:

```hcl
test "my_test" {
  module = "../modules/vpc"
  
  assert "vpc_cidr_matches_input" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
}
```

**Assertion Components**:
- **Name**: Descriptive assertion name
- **Actual**: What to test (output, variable, resource attribute)
- **Expected**: Expected value or reference
- **Condition**: How to compare actual vs expected

## 🔄 Testing Modes

### Mock Mode (Default)
- **Purpose**: Fast unit testing
- **Execution**: Injects mocks, runs `terraform plan`
- **Speed**: Very fast
- **Dependencies**: None (except Terraform)
- **Use Case**: Testing module logic, variable handling

### Live Mode
- **Purpose**: Integration testing
- **Execution**: Runs `terraform init`, `plan`, `apply`
- **Speed**: Slower
- **Dependencies**: Cloud provider credentials
- **Use Case**: End-to-end testing, provider compatibility

## 📊 Reference Types

### Output References
Access module outputs:

```hcl
actual = "output.output_name"
```

**Examples**:
- `"output.vpc_id"`
- `"output.subnet_ids"`
- `"output.security_group_rules"`

### Variable References
Access module variables:

```hcl
expected = "var.variable_name"
```

**Examples**:
- `"var.cidr_block"`
- `"var.instance_count"`
- `"var.tags"`

### Resource References (Live Mode)
Access resource attributes:

```hcl
actual = "resource.aws_vpc.main.id"
```

**Examples**:
- `"resource.aws_vpc.main.cidr_block"`
- `"resource.aws_instance.web.private_ip"`
- `"resource.aws_security_group.web.tags"`

### Literal Values
Direct values:

```hcl
expected = "10.0.0.0/16"
expected = 42
expected = true
expected = ["item1", "item2"]
expected = { key = "value" }
```

## 🎯 Assertion Conditions

### String Conditions
- **`equals`**: Exact string match
- **`contains`**: Substring match
- **`matches`**: Regex pattern match

### JSON Conditions
- **`json_equals`**: Deep JSON comparison

### Examples

```hcl
# String equality
assert "exact_match" {
  actual    = "output.vpc_id"
  expected  = "vpc-12345678"
  condition = "equals"
}

# Substring match
assert "contains_subnet" {
  actual    = "output.subnet_ids"
  expected  = "subnet-"
  condition = "contains"
}

# Regex match
assert "valid_vpc_id" {
  actual    = "output.vpc_id"
  expected  = "vpc-[0-9a-f]{8}"
  condition = "matches"
}

# JSON comparison
assert "tags_match" {
  actual    = "output.tags"
  expected  = '{"Environment":"test","Owner":"team"}'
  condition = "json_equals"
}
```

## 🏃 Execution Flow

### Mock Mode Flow
1. **Parse Test**: Read and parse `.tfunittest.hcl` files
2. **Create Sandbox**: Create temporary directory
3. **Copy Module**: Copy module files to sandbox
4. **Inject Mocks**: Generate mock resource files
5. **Inject Variables**: Create `.tfvars.json` file
6. **Run Terraform**: Execute `terraform init` and `terraform plan`
7. **Parse Output**: Extract outputs from plan JSON
8. **Evaluate Assertions**: Run assertions against outputs
9. **Report Results**: Generate test results

### Live Mode Flow
1. **Parse Test**: Read and parse `.tfunittest.hcl` files
2. **Create Sandbox**: Create temporary directory
3. **Copy Module**: Copy module files to sandbox
4. **Inject Variables**: Create `.tfvars.json` file
5. **Run Terraform**: Execute `terraform init`, `plan`, `apply`
6. **Parse Output**: Extract outputs from state
7. **Evaluate Assertions**: Run assertions against outputs
8. **Cleanup**: Destroy resources (optional)
9. **Report Results**: Generate test results

## 🔧 Configuration

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
export TF_IN_AUTOMATION=true

# Cloud provider credentials (for live mode)
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1
```

## 🎨 Best Practices

### Test Organization
- **One test per file**: Keep tests focused
- **Descriptive names**: Use clear, descriptive test names
- **Group related tests**: Keep similar tests together
- **Use comments**: Document complex test scenarios

### Mock Design
- **Realistic values**: Use realistic mock values
- **Complete attributes**: Include all necessary attributes
- **Consistent naming**: Use consistent resource naming

### Assertion Strategy
- **Test outputs**: Focus on module outputs
- **Test edge cases**: Include boundary conditions
- **Test error cases**: Verify error handling
- **Use appropriate conditions**: Choose the right assertion condition

### Performance
- **Use mock mode**: For fast feedback during development
- **Use live mode**: For integration testing
- **Parallel execution**: Tests run in parallel by default
- **Minimize dependencies**: Keep tests focused and fast

## 🚨 Common Pitfalls

### Mock Mode Issues
- **Missing mocks**: Ensure all required resources are mocked
- **Incorrect attributes**: Mock attributes must match expected structure
- **Resource references**: Use correct resource type and name format

### Live Mode Issues
- **Missing credentials**: Ensure cloud provider credentials are set
- **Resource conflicts**: Avoid naming conflicts between tests
- **Cleanup failures**: Ensure proper resource cleanup

### Assertion Issues
- **Wrong reference format**: Use correct reference syntax
- **Type mismatches**: Ensure actual and expected types match
- **Complex comparisons**: Use `json_equals` for complex structures

## 🔍 Debugging

### Enable Logging
```bash
export TF_LOG=DEBUG
infractest -dir tests
```

### Check JSON Output
```bash
infractest -dir tests -json results.json
cat results.json | jq '.'
```

### Verify Module Structure
```bash
# Check if module is valid
cd modules/vpc
terraform init
terraform validate
```

## 📚 Next Steps

Now that you understand the basic concepts:

1. **[Writing Tests](writing-tests.md)** - Learn how to write effective tests
2. **[Test Modes](test-modes.md)** - Understand when to use each mode
3. **[Assertions](assertions.md)** - Master all assertion conditions
4. **[Basic Examples](../examples/basic.md)** - See practical examples
5. **[Troubleshooting](../../troubleshooting/common-issues.md)** - Solve common problems
