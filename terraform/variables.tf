variable "region" {
  type        = string
  description = "Which region we are deploying to"
  default     = "eu-west-2"
}

variable "vpc_cidr_block" {
  type        = string
  description = "Which VPC CIDR will be deployed into"
  default     = "10.0.0.0/16"
}

variable "environment" {
  type        = string
  description = "Which environment will be deployed into"
  default     = "development"
}

variable "service_name" {
  type        = string
  description = "Name of the service being deployed"
  default     = "user-mgmt-service-api"
}

variable "db_engine_version" {
  type        = string
  description = "Database engine version"
  default     = "16.1"
}

variable "db_instance_class" {
  type        = string
  description = "The instance type of the RDS instance"
  default     = "db.t4g.medium"
}

variable "db_database_name" {
  type        = string
  description = "The name of the database to create when the DB instance is created"
  default     = "user_mgmt_db"
}

variable "db_super_user" {
  type        = string
  description = "Super user for the RDS cluster"
  default     = "postgres"
}