# Quick Start Guide

Get up and running with Infractest in just a few minutes! This guide will walk you through creating and running your first test.

## 🚀 Prerequisites

Before starting, ensure you have:

- Infractest installed (see [Installation Guide](../installation/README.md))
- A basic understanding of Terraform
- A text editor

## 📁 Project Structure

Let's start with a simple project structure:

```
my-terraform-project/
├── modules/
│   └── vpc/
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
└── tests/
    └── vpc_test.tfunittest.hcl
```

## 🏗 Step 1: Create a Simple Terraform Module

First, let's create a basic VPC module to test:

**`modules/vpc/main.tf`**:
```hcl
variable "cidr_block" {
  description = "The CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "name" {
  description = "Name of the VPC"
  type        = string
  default     = "example-vpc"
}

resource "aws_vpc" "main" {
  cidr_block           = var.cidr_block
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = var.name
  }
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "vpc_name" {
  description = "Name of the VPC"
  value       = aws_vpc.main.tags.Name
}
```

## 🧪 Step 2: Create Your First Test

Now, let's create a test file:

**`tests/vpc_test.tfunittest.hcl`**:
```hcl
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

  assert "vpc_cidr_matches_input" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }

  assert "vpc_name_matches_input" {
    actual    = "output.vpc_name"
    expected  = "var.name"
    condition = "equals"
  }
}
```

## 🏃 Step 3: Run Your First Test

Now let's run the test:

```bash
# Navigate to your project directory
cd my-terraform-project

# Run the test in mock mode (default)
infractest -dir tests
```

You should see output similar to:

```
Running tests in mock mode...
✓ vpc_test.tfunittest.hcl::vpc_creation_with_defaults
  ✓ vpc_id_is_generated
  ✓ vpc_cidr_matches_input
  ✓ vpc_name_matches_input

Tests completed: 1 passed, 0 failed
```

## 🎯 Understanding the Test

Let's break down what happened:

### Test Structure
- **`test "vpc_creation_with_defaults"`**: Defines a test case with a descriptive name
- **`module = "../modules/vpc"`**: Points to the Terraform module to test
- **`vars = { ... }`**: Passes variables to the module during testing

### Mock Resources
- **`mock "aws_vpc.main"`**: Creates a mock version of the `aws_vpc.main` resource
- **`attributes = { ... }`**: Defines the mock resource's attributes

### Assertions
- **`assert "vpc_id_is_generated"`**: Tests that the VPC ID output matches expected value
- **`actual = "output.vpc_id"`**: References the module's output
- **`expected = "vpc-12345678"`**: The expected value
- **`condition = "equals"`**: The assertion condition to use

## 🔄 Step 4: Add More Test Cases

Let's add more comprehensive tests:

**`tests/vpc_test.tfunittest.hcl`** (updated):
```hcl
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

  assert "vpc_cidr_matches_input" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }

  assert "vpc_name_matches_input" {
    actual    = "output.vpc_name"
    expected  = "var.name"
    condition = "equals"
  }
}

test "vpc_creation_with_custom_cidr" {
  module = "../modules/vpc"

  vars = {
    cidr_block = "172.16.0.0/16"
    name       = "custom-vpc"
  }

  mock "aws_vpc.main" {
    attributes = {
      id         = "vpc-87654321"
      cidr_block = "172.16.0.0/16"
      tags = {
        Name = "custom-vpc"
      }
    }
  }

  assert "vpc_uses_custom_cidr" {
    actual    = "output.vpc_cidr"
    expected  = "172.16.0.0/16"
    condition = "equals"
  }

  assert "vpc_has_custom_name" {
    actual    = "output.vpc_name"
    expected  = "custom-vpc"
    condition = "equals"
  }
}

test "vpc_creation_with_invalid_cidr" {
  module = "../modules/vpc"

  vars = {
    cidr_block = "invalid-cidr"
    name       = "invalid-vpc"
  }

  # This test should fail - we're testing error handling
  assert "vpc_creation_fails_with_invalid_cidr" {
    actual    = "output.vpc_cidr"
    expected  = "invalid-cidr"
    condition = "equals"
  }
}
```

Run the updated tests:

```bash
infractest -dir tests
```

## 📊 Step 5: Generate JSON Report

For CI/CD integration or detailed analysis, generate a JSON report:

```bash
infractest -dir tests -json test-results.json
```

This creates a detailed JSON report with test results, timing, and error details.

## 🎛 Step 6: Try Different Test Modes

### Mock Mode (Default)
```bash
infractest -dir tests -mode mock
```

### Live Mode (Requires AWS credentials)
```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1

# Run in live mode
infractest -dir tests -mode live
```

## 🔍 Step 7: Debug Test Failures

If a test fails, you can debug it:

1. **Check the error output**:
   ```bash
   infractest -dir tests -mode mock
   ```

2. **Enable Terraform logging**:
   ```bash
   export TF_LOG=DEBUG
   infractest -dir tests -mode mock
   ```

3. **Check the JSON report**:
   ```bash
   infractest -dir tests -json results.json
   cat results.json | jq '.'
   ```

## 🎉 Congratulations!

You've successfully:

- ✅ Created your first Terraform module
- ✅ Written comprehensive tests for it
- ✅ Run tests in both mock and live modes
- ✅ Generated detailed reports
- ✅ Debugged test failures

## 🚀 Next Steps

Now that you have the basics down, explore:

1. **[Writing Tests](writing-tests.md)** - Learn advanced testing techniques
2. **[Test Modes](test-modes.md)** - Understand mock vs live testing
3. **[Assertions](assertions.md)** - Master all assertion conditions
4. **[Basic Examples](../examples/basic.md)** - See more examples
5. **[CI/CD Integration](ci-cd.md)** - Automate your testing

## 🆘 Need Help?

If you run into issues:

1. Check the [Troubleshooting Guide](../../troubleshooting/common-issues.md)
2. Review the [FAQ](../../troubleshooting/faq.md)
3. Ask questions in [GitHub Discussions](https://github.com/memetics19/infractest/discussions)
4. Report bugs in [GitHub Issues](https://github.com/memetics19/infractest/issues)

Happy testing! 🧪
