module "ecs_cluster" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 5.0"

  # Override the default of enabling Container Insights to keep costs down whilst this is just for E2E tests
  cluster_settings = []

  cluster_name = "${var.unique_identifier_prefix}-${var.environment}"
}

module "ecs_service" {
  source  = "terraform-aws-modules/ecs/aws//modules/service"
  version = "~> 5.0"

  name        = local.name
  cluster_arn = module.ecs_cluster.cluster_arn

  cpu    = var.fargate_task_cpu
  memory = var.fargate_task_memory

  enable_autoscaling = false

  container_definitions = {
    main = {
      essential       = true
      image           = var.fargate_docker_image
      cpuArchitecture = var.fargate_cpu_architecture
      port_mappings = [
        {
          name          = "http"
          containerPort = var.fargate_container_port
          hostPort      = var.fargate_container_port
          protocol      = "tcp"
        }
      ]
      enable_cloudwatch_logging = true

      environment = [
        {
          name : "database_host_name"
          value : module.db.cluster_endpoint
        },
        {
          name : "database_port"
          value : "5432"
        },
        {
          name : "database_name"
          value : var.db_database_name
        },
        {
          name : "database_username"
          value : var.db_master_username
        },
        {
          name : "database_ssl_mode"
          value : "disable"
        }
      ]

      # Env vars mounted from Secrets Manager or SSM Parameter Store
      # ECS does not support querying JSON structures natively, so storing a single password value to keep it simple
      "secrets" : [
        {
          "name" : "database_password",
          "valueFrom" : aws_secretsmanager_secret_version.db_master_password.arn
        }
      ]
    }
  }

  load_balancer = {
    service = {
      target_group_arn = module.alb.target_groups["ecs"].arn
      container_name   = "main"
      container_port   = var.fargate_container_port
    }
  }

  subnet_ids = module.vpc.private_subnets

  security_group_name            = "${local.name}-fargate-task"
  security_group_description     = "Assigned to ${local.name} Fargate task"
  security_group_use_name_prefix = false
  security_group_rules = {
    alb_ingress_8080 = {
      type                     = "ingress"
      from_port                = var.fargate_container_port
      to_port                  = var.fargate_container_port
      protocol                 = "tcp"
      description              = "Application traffic"
      source_security_group_id = module.alb.security_group_id
    }
    egress_all = {
      type        = "egress"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }
}