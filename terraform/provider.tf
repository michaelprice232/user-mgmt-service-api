provider "aws" {
  region = var.region

  default_tags {
    tags = {
      application = var.service_name
      environment = var.environment
      owner       = "Michael Price"
    }
  }
}

# Do not use any locking to enable parallel runs of the E2E tests in the CI system
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
  required_version = "~> 1.10.0"
}