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
  default     = "dev"
}

variable "service_name" {
  type        = string
  description = "Name of the service being deployed"
  default     = "user-mgmt-service-api"
}

variable "unique_identifier_prefix" {
  type        = string
  description = "Unique identifier prefix to allow parallel deployments into AWS via the CI system without resource name clashes"
}

variable "db_engine_version" {
  type        = string
  description = "Database engine version"
  default     = "16.6"
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

variable "db_master_username" {
  type        = string
  description = "Master username for the RDS cluster"
  default     = "postgres"
}

variable "fargate_task_cpu" {
  type        = number
  description = "Amount of CPU to allocate to the Fargate task"
  default     = 256 # 0.25 core
}

variable "fargate_task_memory" {
  type        = number
  description = "Amount of memory to allocate to the Fargate task"
  default     = 512 # 0.5 GB
}

variable "fargate_docker_image" {
  type        = string
  description = "Docker image URL to run"
  default     = "633681147894.dkr.ecr.eu-west-2.amazonaws.com/user-mgmt-service-api:108c7d2ae11204458fce0341c740e07898959ffe"
}

variable "fargate_container_port" {
  type        = number
  description = "The port number that the app runs on"
  default     = 8080
}

variable "fargate_cpu_architecture" {
  type        = string
  description = "CPU architecture type to run with the Fargate task"
  default     = "ARM64"
}