# Writing Tests

This guide covers how to write effective and comprehensive tests for your Terraform modules using Infractest.

## 📝 Test File Structure

### Basic Test File

```hcl
# File: tests/vpc_test.tfunittest.hcl

test "vpc_creation_with_defaults" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
    name       = "test-vpc"
  }
  
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-12345678"
      cidr_block = "10.0.0.0/16"
      tags = {
        Name = "test-vpc"
      }
    }
  }
  
  assert "vpc_id_is_generated" {
    actual    = "output.vpc_id"
    expected  = "vpc-12345678"
    condition = "equals"
  }
}
```

### Multiple Tests in One File

```hcl
# File: tests/vpc_comprehensive.tfunittest.hcl

test "vpc_creation_with_defaults" {
  # Test implementation
}

test "vpc_creation_with_custom_cidr" {
  # Test implementation
}

test "vpc_creation_with_tags" {
  # Test implementation
}
```

## 🎯 Test Naming Conventions

### Descriptive Names
Use clear, descriptive names that explain what the test validates:

```hcl
# Good
test "vpc_creation_with_defaults"
test "vpc_creation_with_custom_cidr"
test "vpc_creation_fails_with_invalid_cidr"
test "vpc_has_correct_tags"
test "vpc_enables_dns_support"

# Avoid
test "test1"
test "vpc_test"
test "basic_test"
```

### Naming Patterns
- **Feature + Condition**: `vpc_creation_with_defaults`
- **Feature + Error**: `vpc_creation_fails_with_invalid_cidr`
- **Resource + Attribute**: `vpc_has_correct_tags`
- **Behavior + State**: `vpc_enables_dns_support`

## 🏗 Module References

### Relative Paths
```hcl
test "my_test" {
  module = "../modules/vpc"           # Go up one level
  module = "../../modules/vpc"        # Go up two levels
  module = "./local-module"           # Current directory
}
```

### Absolute Paths
```hcl
test "my_test" {
  module = "/full/path/to/module"     # Absolute path
}
```

### Module Validation
Ensure your module is valid before testing:

```bash
# Validate the module
cd modules/vpc
terraform init
terraform validate
```

## 📊 Variable Configuration

### Basic Variables
```hcl
test "my_test" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
    name       = "test-vpc"
    enable_dns = true
  }
}
```

### Complex Variables
```hcl
test "my_test" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
    name       = "test-vpc"
    
    # List variables
    availability_zones = ["us-east-1a", "us-east-1b"]
    
    # Map variables
    tags = {
      Environment = "test"
      Owner       = "team"
      Project     = "infrastructure"
    }
    
    # Nested maps
    subnet_config = {
      public = {
        cidr = "10.0.1.0/24"
        az   = "us-east-1a"
      }
      private = {
        cidr = "10.0.2.0/24"
        az   = "us-east-1b"
      }
    }
  }
}
```

### Variable Types
```hcl
vars = {
  # String
  name = "test-vpc"
  
  # Number
  instance_count = 3
  port_number    = 8080
  
  # Boolean
  enable_logging = true
  create_nat     = false
  
  # List
  subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  
  # Map
  tags = {
    Environment = "test"
    Owner       = "team"
  }
}
```

## 🎭 Mock Resources

### Basic Mock Structure
```hcl
mock "resource_type.name" {
  attributes = {
    attribute_name = "value"
    another_attr   = "another_value"
  }
}
```

### Common Resource Mocks

#### VPC Mock
```hcl
mock "aws_vpc.main" {
  attributes = {
    id                = "vpc-12345678"
    cidr_block        = "10.0.0.0/16"
    enable_dns_hostnames = true
    enable_dns_support   = true
    tags = {
      Name        = "test-vpc"
      Environment = "test"
    }
  }
}
```

#### Subnet Mock
```hcl
mock "aws_subnet.public" {
  attributes = {
    id               = "subnet-12345678"
    vpc_id           = "vpc-12345678"
    cidr_block       = "10.0.1.0/24"
    availability_zone = "us-east-1a"
    map_public_ip_on_launch = true
    tags = {
      Name = "public-subnet"
      Type = "public"
    }
  }
}
```

#### Security Group Mock
```hcl
mock "aws_security_group.web" {
  attributes = {
    id          = "sg-12345678"
    name        = "web-sg"
    description = "Security group for web servers"
    vpc_id      = "vpc-12345678"
    tags = {
      Name = "web-security-group"
    }
  }
}
```

#### EC2 Instance Mock
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
      Name = "web-server"
      Role = "web"
    }
  }
}
```

### Mock Best Practices

#### Use Realistic Values
```hcl
# Good - realistic values
mock "aws_vpc.main" {
  attributes = {
    id         = "vpc-12345678"        # Realistic VPC ID format
    cidr_block = "10.0.0.0/16"         # Valid CIDR block
  }
}

# Avoid - unrealistic values
mock "aws_vpc.main" {
  attributes = {
    id         = "vpc-123"             # Too short
    cidr_block = "invalid-cidr"        # Invalid CIDR
  }
}
```

#### Include All Required Attributes
```hcl
# Good - complete mock
mock "aws_vpc.main" {
  attributes = {
    id                = "vpc-12345678"
    cidr_block        = "10.0.0.0/16"
    enable_dns_hostnames = true
    enable_dns_support   = true
    tags = {
      Name = "test-vpc"
    }
  }
}

# Avoid - incomplete mock
mock "aws_vpc.main" {
  attributes = {
    id = "vpc-12345678"
    # Missing other attributes
  }
}
```

## ✅ Assertions

### Basic Assertion Structure
```hcl
assert "descriptive_name" {
  actual    = "output.output_name"
  expected  = "expected_value"
  condition = "equals"
}
```

### Output Assertions
```hcl
# Test module outputs
assert "vpc_id_is_generated" {
  actual    = "output.vpc_id"
  expected  = "vpc-12345678"
  condition = "equals"
}

assert "vpc_cidr_matches_input" {
  actual    = "output.vpc_cidr"
  expected  = "var.cidr_block"
  condition = "equals"
}
```

### Variable Assertions
```hcl
# Test that outputs match input variables
assert "vpc_name_matches_input" {
  actual    = "output.vpc_name"
  expected  = "var.name"
  condition = "equals"
}
```

### Resource Assertions (Live Mode)
```hcl
# Test resource attributes in live mode
assert "vpc_has_correct_cidr" {
  actual    = "resource.aws_vpc.main.cidr_block"
  expected  = "10.0.0.0/16"
  condition = "equals"
}
```

### Complex Assertions
```hcl
# Test complex outputs
assert "vpc_tags_are_correct" {
  actual    = "output.vpc_tags"
  expected  = '{"Name":"test-vpc","Environment":"test"}'
  condition = "json_equals"
}

# Test list outputs
assert "subnet_ids_are_generated" {
  actual    = "output.subnet_ids"
  expected  = "subnet-"
  condition = "contains"
}
```

## 🎨 Test Organization

### Single Responsibility
Each test should focus on one specific behavior:

```hcl
# Good - focused test
test "vpc_creation_with_defaults" {
  # Tests only default VPC creation
}

test "vpc_creation_with_custom_cidr" {
  # Tests only custom CIDR handling
}

# Avoid - testing multiple things
test "vpc_creation_and_subnets" {
  # Tests both VPC and subnet creation
}
```

### Test Groups
Group related tests together:

```hcl
# File: tests/vpc_basic.tfunittest.hcl
test "vpc_creation_with_defaults" { }
test "vpc_creation_with_custom_cidr" { }
test "vpc_creation_with_tags" { }

# File: tests/vpc_advanced.tfunittest.hcl
test "vpc_with_nat_gateway" { }
test "vpc_with_vpc_endpoints" { }
test "vpc_with_flow_logs" { }
```

### Test Dependencies
Avoid test dependencies - each test should be independent:

```hcl
# Good - independent tests
test "vpc_creation" {
  # Creates its own VPC
}

test "subnet_creation" {
  # Creates its own VPC and subnet
}

# Avoid - dependent tests
test "vpc_creation" {
  # Creates VPC
}

test "subnet_creation" {
  # Depends on VPC from previous test
}
```

## 🔍 Edge Cases and Error Testing

### Valid Input Testing
```hcl
test "vpc_creation_with_valid_cidr" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
  }
  
  # Test that valid input works
  assert "vpc_created_successfully" {
    actual    = "output.vpc_id"
    expected  = "vpc-12345678"
    condition = "equals"
  }
}
```

### Invalid Input Testing
```hcl
test "vpc_creation_with_invalid_cidr" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "invalid-cidr"
  }
  
  # Test that invalid input is handled
  assert "vpc_creation_fails" {
    actual    = "output.error_message"
    expected  = "invalid CIDR"
    condition = "contains"
  }
}
```

### Boundary Testing
```hcl
test "vpc_creation_with_minimum_cidr" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/28"  # Minimum valid CIDR
  }
  
  assert "vpc_created_with_minimum_cidr" {
    actual    = "output.vpc_cidr"
    expected  = "10.0.0.0/28"
    condition = "equals"
  }
}
```

## 📝 Comments and Documentation

### Test Documentation
```hcl
# Test VPC creation with various CIDR blocks
# This test validates that the VPC module correctly handles
# different CIDR block inputs and generates appropriate outputs

test "vpc_creation_with_custom_cidr" {
  module = "../modules/vpc"
  
  # Test with a non-default CIDR block
  vars = {
    cidr_block = "172.16.0.0/16"
    name       = "custom-vpc"
  }
  
  # Mock the VPC resource with expected attributes
  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-87654321"
      cidr_block = "172.16.0.0/16"
      tags = {
        Name = "custom-vpc"
      }
    }
  }
  
  # Verify the VPC ID is generated
  assert "vpc_id_is_generated" {
    actual    = "output.vpc_id"
    expected  = "vpc-87654321"
    condition = "equals"
  }
  
  # Verify the CIDR block matches input
  assert "vpc_cidr_matches_input" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
}
```

## 🚀 Performance Considerations

### Mock Mode for Development
```hcl
# Use mock mode for fast feedback during development
test "vpc_creation" {
  module = "../modules/vpc"
  # ... test configuration
}
```

### Live Mode for Integration
```hcl
# Use live mode for integration testing
test "vpc_creation_integration" {
  module = "../modules/vpc"
  # ... test configuration
}
```

### Parallel Execution
Tests run in parallel by default - design tests to be independent:

```hcl
# These tests can run in parallel
test "vpc_creation_test1" { }
test "vpc_creation_test2" { }
test "vpc_creation_test3" { }
```

## 🔧 Debugging Tests

### Enable Terraform Logging
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
cd modules/vpc
terraform init
terraform validate
terraform plan
```

## 📚 Best Practices Summary

1. **Use descriptive test names** that explain what's being tested
2. **Keep tests focused** - one behavior per test
3. **Use realistic mock values** that match real resource attributes
4. **Test both success and failure cases**
5. **Include comprehensive assertions** for all important outputs
6. **Document complex test scenarios** with comments
7. **Organize tests logically** by feature or behavior
8. **Make tests independent** - no dependencies between tests
9. **Use appropriate test modes** - mock for development, live for integration
10. **Validate your modules** before writing tests

## 🎯 Next Steps

Now that you know how to write tests:

1. **[Test Modes](test-modes.md)** - Understand when to use each mode
2. **[Assertions](assertions.md)** - Master all assertion conditions
3. **[Basic Examples](../examples/basic.md)** - See practical examples
4. **[Advanced Examples](../examples/advanced.md)** - Learn complex scenarios
5. **[Troubleshooting](../../troubleshooting/common-issues.md)** - Solve common problems
