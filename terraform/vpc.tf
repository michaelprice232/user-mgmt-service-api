module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name               = "${var.unique_identifier_prefix}-${var.environment}"
  cidr               = var.vpc_cidr_block
  azs                = [for i in range(3) : data.aws_availability_zones.available.names[i]]
  private_subnets    = [for i in range(3) : cidrsubnet(var.vpc_cidr_block, 4, i)]
  public_subnets     = [for i in range(3, 6) : cidrsubnet(var.vpc_cidr_block, 4, i)]
  enable_nat_gateway = true
  single_nat_gateway = true
}