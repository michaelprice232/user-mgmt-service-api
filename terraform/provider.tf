provider "aws" {
  region = var.region
  # todo: remove for OIDC
  profile = "personal"

  default_tags {
    tags = {
      application = var.service_name
      environment = var.environment
      owner       = "Michael Price"
    }
  }
}