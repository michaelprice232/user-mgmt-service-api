data "aws_availability_zones" "available" {
  state = "available"
}

resource "random_password" "aurora_master_password" {
  length  = 50
  special = false
}