module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "5.17.0"
  name = "main"
  cidr = "10.0.0.0/16"
  azs             = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  enable_nat_gateway = true
  single_nat_gateway = true
}
#
# module "db" {
#   source = "terraform-aws-modules/rds/aws"
#   version = "6.10.0"
#
#   identifier = "demodb"
#
#   engine            = "mysql"
#   engine_version    = "5.7"
#   instance_class    = "db.t3a.large"
#   allocated_storage = 5
#
#   db_name  = "demodb"
#   username = "user"
#   port     = "3306"
#
#   iam_database_authentication_enabled = true
#
#   vpc_security_group_ids = ["sg-12345678"]
#
#   maintenance_window = "Mon:00:00-Mon:03:00"
#   backup_window      = "03:00-06:00"
#
#   # Enhanced Monitoring - see example for details on how to create the role
#   # by yourself, in case you don't want to create it automatically
#   monitoring_interval    = "30"
#   monitoring_role_name   = "MyRDSMonitoringRole"
#   create_monitoring_role = true
#
#   tags = {
#     Owner       = "user"
#     Environment = "dev"
#   }
#
#   # DB subnet group
#   create_db_subnet_group = true
#   subnet_ids             = ["subnet-12345678", "subnet-87654321"]
#
#   # DB parameter group
#   family = "mysql5.7"
#
#   # DB option group
#   major_engine_version = "5.7"
#
#   # Database Deletion Protection
#   deletion_protection = true
#
#   parameters = [
#     {
#       name  = "character_set_client"
#       value = "utf8mb4"
#     },
#     {
#       name  = "character_set_server"
#       value = "utf8mb4"
#     }
#   ]
#
#   options = [
#     {
#       option_name = "MARIADB_AUDIT_PLUGIN"
#
#       option_settings = [
#         {
#           name  = "SERVER_AUDIT_EVENTS"
#           value = "CONNECT"
#         },
#         {
#           name  = "SERVER_AUDIT_FILE_ROTATIONS"
#           value = "37"
#         },
#       ]
#     },
#   ]
# }