provider "aws" {
  region                      = "us-east-1"
  access_key                  = "mock"
  secret_key                  = "mock"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  skip_region_validation      = true
#   s3_force_path_style         = true
  endpoints {
    ec2 = "http://127.0.0.1:4566"
  }
}

variable "mock_mode" {
  type    = bool
  default = false
}

locals {
  create_vpc = var.mock_mode ? false : true
}

resource "aws_vpc" "main" {
  count      = local.create_vpc ? 1 : 0
  cidr_block = var.cidr_block
}

output "vpc_cidr" {
  value = var.mock_mode ? var.cidr_block : aws_vpc.main[0].cidr_block
}

variable "cidr_block" {
  type        = string
  description = "The CIDR block for the VPC"
  default     = "10.0.0.0/16"
}