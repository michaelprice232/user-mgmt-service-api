provider "aws" {
  region  = var.region
  profile = "personal"

  default_tags {
    tags = {
      environment = var.environment
      owner       = "Michael Price"
    }
  }
}