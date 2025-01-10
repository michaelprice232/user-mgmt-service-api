module "db" {
  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "9.11.0"

  name           = local.name
  engine         = "aurora-postgresql"
  engine_version = var.db_engine_version
  instance_class = var.db_instance_class

  instances = {
    one = {} # Single writer instance
  }

  master_username             = var.db_super_user
  manage_master_user_password = false
  master_password             = random_password.aurora_master_password.result

  database_name = var.db_database_name

  vpc_id                 = module.vpc.vpc_id
  create_db_subnet_group = true
  subnets                = module.vpc.private_subnets
  db_subnet_group_name   = local.name

  security_group_name        = "${local.name}-db"
  security_group_description = "Assigned to ${local.name} Aurora cluster"
  security_group_rules = {
    vpc_ingress = {
      cidr_blocks = [var.vpc_cidr_block]
    }
  }
  storage_encrypted = true
  apply_immediately = true
}