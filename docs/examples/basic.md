# Basic Examples

This section provides practical examples of using Infractest for common Terraform testing scenarios.

## 🏗 VPC Module Testing

### Simple VPC Test

**Module**: `modules/vpc/main.tf`
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

**Test**: `tests/vpc_basic.tfunittest.hcl`
```hcl
test "vpc_creation_with_defaults" {
  module = "../modules/vpc"
  
  vars = {
    cidr_block = "10.0.0.0/16"
    name       = "test-vpc"
  }
  
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
```

## 🔒 Security Group Module Testing

### Security Group with Rules

**Module**: `modules/security-group/main.tf`
```hcl
variable "name" {
  description = "Name of the security group"
  type        = string
}

variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "allow_ssh" {
  description = "Whether to allow SSH access"
  type        = bool
  default     = false
}

variable "allow_http" {
  description = "Whether to allow HTTP access"
  type        = bool
  default     = false
}

resource "aws_security_group" "main" {
  name        = var.name
  description = "Security group for ${var.name}"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = var.allow_ssh ? [1] : []
    content {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  dynamic "ingress" {
    for_each = var.allow_http ? [1] : []
    content {
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = var.name
  }
}

output "security_group_id" {
  description = "ID of the security group"
  value       = aws_security_group.main.id
}

output "security_group_name" {
  description = "Name of the security group"
  value       = aws_security_group.main.name
}

output "ingress_rules" {
  description = "Ingress rules"
  value       = aws_security_group.main.ingress
}
```

**Test**: `tests/security_group_basic.tfunittest.hcl`
```hcl
test "security_group_creation_with_ssh" {
  module = "../modules/security-group"
  
  vars = {
    name      = "web-sg"
    vpc_id    = "vpc-12345678"
    allow_ssh = true
    allow_http = false
  }
  
  mock "aws_security_group.main" {
    attributes = {
      id          = "sg-12345678"
      name        = "web-sg"
      description = "Security group for web-sg"
      vpc_id      = "vpc-12345678"
      ingress = [
        {
          from_port   = 22
          to_port     = 22
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
        Name = "web-sg"
      }
    }
  }
  
  assert "security_group_id_is_generated" {
    actual    = "output.security_group_id"
    expected  = "sg-12345678"
    condition = "equals"
  }
  
  assert "security_group_has_ssh_rule" {
    actual    = "output.ingress_rules"
    expected  = "22"
    condition = "contains"
  }
  
  assert "security_group_name_matches_input" {
    actual    = "output.security_group_name"
    expected  = "var.name"
    condition = "equals"
  }
}

test "security_group_creation_with_http" {
  module = "../modules/security-group"
  
  vars = {
    name      = "web-sg"
    vpc_id    = "vpc-12345678"
    allow_ssh = false
    allow_http = true
  }
  
  mock "aws_security_group.main" {
    attributes = {
      id          = "sg-87654321"
      name        = "web-sg"
      description = "Security group for web-sg"
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
        Name = "web-sg"
      }
    }
  }
  
  assert "security_group_has_http_rule" {
    actual    = "output.ingress_rules"
    expected  = "80"
    condition = "contains"
  }
  
  assert "security_group_no_ssh_rule" {
    actual    = "output.ingress_rules"
    expected  = "22"
    condition = "contains"
  }
}
```

## 🖥 EC2 Instance Module Testing

### EC2 Instance with User Data

**Module**: `modules/ec2-instance/main.tf`
```hcl
variable "ami_id" {
  description = "AMI ID for the instance"
  type        = string
}

variable "instance_type" {
  description = "Instance type"
  type        = string
  default     = "t3.micro"
}

variable "subnet_id" {
  description = "Subnet ID for the instance"
  type        = string
}

variable "security_group_ids" {
  description = "Security group IDs"
  type        = list(string)
}

variable "user_data" {
  description = "User data script"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags for the instance"
  type        = map(string)
  default     = {}
}

resource "aws_instance" "main" {
  ami                    = var.ami_id
  instance_type          = var.instance_type
  subnet_id              = var.subnet_id
  vpc_security_group_ids = var.security_group_ids
  user_data              = var.user_data

  tags = merge(var.tags, {
    Name = "ec2-instance"
  })
}

output "instance_id" {
  description = "ID of the instance"
  value       = aws_instance.main.id
}

output "instance_public_ip" {
  description = "Public IP of the instance"
  value       = aws_instance.main.public_ip
}

output "instance_private_ip" {
  description = "Private IP of the instance"
  value       = aws_instance.main.private_ip
}

output "instance_tags" {
  description = "Tags of the instance"
  value       = aws_instance.main.tags
}
```

**Test**: `tests/ec2_instance_basic.tfunittest.hcl`
```hcl
test "ec2_instance_creation_with_defaults" {
  module = "../modules/ec2-instance"
  
  vars = {
    ami_id             = "ami-12345678"
    instance_type      = "t3.micro"
    subnet_id          = "subnet-12345678"
    security_group_ids = ["sg-12345678"]
    tags = {
      Environment = "test"
      Owner       = "team"
    }
  }
  
  mock "aws_instance.main" {
    attributes = {
      id               = "i-12345678"
      ami              = "ami-12345678"
      instance_type    = "t3.micro"
      subnet_id        = "subnet-12345678"
      vpc_security_group_ids = ["sg-12345678"]
      public_ip        = "203.0.113.100"
      private_ip       = "10.0.1.100"
      tags = {
        Name        = "ec2-instance"
        Environment = "test"
        Owner       = "team"
      }
    }
  }
  
  assert "instance_id_is_generated" {
    actual    = "output.instance_id"
    expected  = "i-12345678"
    condition = "equals"
  }
  
  assert "instance_has_public_ip" {
    actual    = "output.instance_public_ip"
    expected  = "203.0.113.100"
    condition = "equals"
  }
  
  assert "instance_has_private_ip" {
    actual    = "output.instance_private_ip"
    expected  = "10.0.1.100"
    condition = "equals"
  }
  
  assert "instance_tags_include_environment" {
    actual    = "output.instance_tags"
    expected  = "Environment"
    condition = "contains"
  }
}

test "ec2_instance_creation_with_user_data" {
  module = "../modules/ec2-instance"
  
  vars = {
    ami_id             = "ami-12345678"
    instance_type      = "t3.small"
    subnet_id          = "subnet-12345678"
    security_group_ids = ["sg-12345678"]
    user_data          = "#!/bin/bash\necho 'Hello World'"
    tags = {
      Environment = "production"
      Owner       = "team"
    }
  }
  
  mock "aws_instance.main" {
    attributes = {
      id               = "i-87654321"
      ami              = "ami-12345678"
      instance_type    = "t3.small"
      subnet_id        = "subnet-12345678"
      vpc_security_group_ids = ["sg-12345678"]
      user_data        = "#!/bin/bash\necho 'Hello World'"
      public_ip        = "203.0.113.200"
      private_ip       = "10.0.1.200"
      tags = {
        Name        = "ec2-instance"
        Environment = "production"
        Owner       = "team"
      }
    }
  }
  
  assert "instance_has_user_data" {
    actual    = "output.instance_tags"
    expected  = "production"
    condition = "contains"
  }
  
  assert "instance_type_matches_input" {
    actual    = "output.instance_tags"
    expected  = "t3.small"
    condition = "contains"
  }
}
```

## 🗄 RDS Database Module Testing

### RDS Instance with Subnet Group

**Module**: `modules/rds-instance/main.tf`
```hcl
variable "identifier" {
  description = "Identifier for the RDS instance"
  type        = string
}

variable "engine" {
  description = "Database engine"
  type        = string
  default     = "mysql"
}

variable "engine_version" {
  description = "Database engine version"
  type        = string
  default     = "8.0"
}

variable "instance_class" {
  description = "Instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Name of the database"
  type        = string
}

variable "username" {
  description = "Master username"
  type        = string
}

variable "password" {
  description = "Master password"
  type        = string
  sensitive   = true
}

variable "vpc_security_group_ids" {
  description = "VPC security group IDs"
  type        = list(string)
}

variable "db_subnet_group_name" {
  description = "DB subnet group name"
  type        = string
}

variable "backup_retention_period" {
  description = "Backup retention period in days"
  type        = number
  default     = 7
}

resource "aws_db_instance" "main" {
  identifier             = var.identifier
  engine                 = var.engine
  engine_version         = var.engine_version
  instance_class         = var.instance_class
  allocated_storage      = var.allocated_storage
  db_name                = var.db_name
  username               = var.username
  password               = var.password
  vpc_security_group_ids = var.vpc_security_group_ids
  db_subnet_group_name   = var.db_subnet_group_name
  backup_retention_period = var.backup_retention_period
  skip_final_snapshot    = true

  tags = {
    Name = var.identifier
  }
}

output "db_instance_id" {
  description = "ID of the DB instance"
  value       = aws_db_instance.main.id
}

output "db_instance_endpoint" {
  description = "Endpoint of the DB instance"
  value       = aws_db_instance.main.endpoint
}

output "db_instance_arn" {
  description = "ARN of the DB instance"
  value       = aws_db_instance.main.arn
}

output "db_instance_status" {
  description = "Status of the DB instance"
  value       = aws_db_instance.main.status
}
```

**Test**: `tests/rds_instance_basic.tfunittest.hcl`
```hcl
test "rds_instance_creation_with_mysql" {
  module = "../modules/rds-instance"
  
  vars = {
    identifier             = "test-db"
    engine                = "mysql"
    engine_version        = "8.0"
    instance_class        = "db.t3.micro"
    allocated_storage     = 20
    db_name               = "testdb"
    username              = "admin"
    password              = "password123"
    vpc_security_group_ids = ["sg-12345678"]
    db_subnet_group_name  = "test-db-subnet-group"
    backup_retention_period = 7
  }
  
  mock "aws_db_instance.main" {
    attributes = {
      id                   = "test-db"
      identifier           = "test-db"
      engine               = "mysql"
      engine_version       = "8.0"
      instance_class       = "db.t3.micro"
      allocated_storage    = 20
      db_name              = "testdb"
      username             = "admin"
      vpc_security_group_ids = ["sg-12345678"]
      db_subnet_group_name = "test-db-subnet-group"
      backup_retention_period = 7
      endpoint             = "test-db.123456789012.us-east-1.rds.amazonaws.com"
      arn                  = "arn:aws:rds:us-east-1:123456789012:db:test-db"
      status               = "available"
      tags = {
        Name = "test-db"
      }
    }
  }
  
  assert "db_instance_id_matches_input" {
    actual    = "output.db_instance_id"
    expected  = "test-db"
    condition = "equals"
  }
  
  assert "db_instance_has_endpoint" {
    actual    = "output.db_instance_endpoint"
    expected  = "test-db.123456789012.us-east-1.rds.amazonaws.com"
    condition = "equals"
  }
  
  assert "db_instance_has_arn" {
    actual    = "output.db_instance_arn"
    expected  = "arn:aws:rds:us-east-1:123456789012:db:test-db"
    condition = "equals"
  }
  
  assert "db_instance_status_is_available" {
    actual    = "output.db_instance_status"
    expected  = "available"
    condition = "equals"
  }
}

test "rds_instance_creation_with_postgresql" {
  module = "../modules/rds-instance"
  
  vars = {
    identifier             = "test-pg-db"
    engine                = "postgres"
    engine_version        = "13.7"
    instance_class        = "db.t3.small"
    allocated_storage     = 50
    db_name               = "testpgdb"
    username              = "postgres"
    password              = "password123"
    vpc_security_group_ids = ["sg-87654321"]
    db_subnet_group_name  = "test-pg-db-subnet-group"
    backup_retention_period = 14
  }
  
  mock "aws_db_instance.main" {
    attributes = {
      id                   = "test-pg-db"
      identifier           = "test-pg-db"
      engine               = "postgres"
      engine_version       = "13.7"
      instance_class       = "db.t3.small"
      allocated_storage    = 50
      db_name              = "testpgdb"
      username             = "postgres"
      vpc_security_group_ids = ["sg-87654321"]
      db_subnet_group_name = "test-pg-db-subnet-group"
      backup_retention_period = 14
      endpoint             = "test-pg-db.123456789012.us-east-1.rds.amazonaws.com"
      arn                  = "arn:aws:rds:us-east-1:123456789012:db:test-pg-db"
      status               = "available"
      tags = {
        Name = "test-pg-db"
      }
    }
  }
  
  assert "db_instance_uses_postgresql" {
    actual    = "output.db_instance_id"
    expected  = "test-pg-db"
    condition = "equals"
  }
  
  assert "db_instance_has_larger_storage" {
    actual    = "output.db_instance_endpoint"
    expected  = "test-pg-db.123456789012.us-east-1.rds.amazonaws.com"
    condition = "equals"
  }
}
```

## 🚀 Running the Examples

### Run All Tests
```bash
# Run all tests in mock mode
infractest -dir tests

# Run all tests in live mode
infractest -dir tests -mode live

# Generate JSON report
infractest -dir tests -json results.json
```

### Run Specific Test Files
```bash
# Run only VPC tests
infractest -dir tests -mode mock

# Run only security group tests
infractest -dir tests -mode mock
```

### Debug Tests
```bash
# Enable debug logging
export TF_LOG=DEBUG
infractest -dir tests -mode mock

# Check JSON output
infractest -dir tests -json results.json
cat results.json | jq '.'
```

## 🎯 Next Steps

These basic examples demonstrate:

1. **Module Testing**: Testing individual Terraform modules
2. **Mock Resources**: Using mocks to simulate cloud resources
3. **Assertions**: Validating module outputs and behavior
4. **Variable Testing**: Testing with different input values
5. **Output Validation**: Ensuring outputs match expectations

To learn more:

1. **[Advanced Examples](advanced.md)** - Complex testing scenarios
2. **[Test Modes](test-modes.md)** - Mock vs live testing
3. **[Assertions](assertions.md)** - All assertion conditions
4. **[Writing Tests](writing-tests.md)** - Best practices
5. **[Troubleshooting](../../troubleshooting/common-issues.md)** - Common problems
