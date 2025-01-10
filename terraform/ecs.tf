module "ecs_cluster" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 5.0"

  cluster_name = local.name
}

module "ecs_service" {
  source  = "terraform-aws-modules/ecs/aws//modules/service"
  version = "~> 5.0"

  name        = local.name
  cluster_arn = module.ecs_cluster.cluster_arn

  cpu    = 1024
  memory = 4096

  enable_autoscaling = false

  # Container definition(s)
  container_definitions = {

    main = {
      # cpu       = 512
      # memory    = 1024
      essential = true
      image     = "633681147894.dkr.ecr.eu-west-2.amazonaws.com/user-mgmt-service-api:817d05ba07adf9aa1e0fd29e4bba8313595ed370"
      port_mappings = [
        {
          name          = "http"
          containerPort = 8080
          hostPort      = 8080
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
          value : var.db_super_user
        },
        // todo: pull from secrets manager?
        {
          name : "database_password"
          value : random_password.aurora_master_password.result
        },
        {
          name : "database_ssl_mode"
          value : "disable"
        }
      ]
    }
  }

  load_balancer = {
    service = {
      target_group_arn = module.alb.target_groups["ecs"].arn
      container_name   = "main"
      container_port   = 8080
    }
  }

  subnet_ids = module.vpc.private_subnets

  security_group_name        = "${local.name}-fargate-task"
  security_group_description = "Assigned to ${local.name} Fargate task"
  security_group_rules = {
    alb_ingress_8080 = {
      type                     = "ingress"
      from_port                = 8080
      to_port                  = 8080
      protocol                 = "tcp"
      description              = "Service port"
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

  service_tags = {
    "ServiceTag" = "Tag on service level"
  }

  # tags = local.tags
}