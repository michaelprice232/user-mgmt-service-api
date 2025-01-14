module "ecs_cluster" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 5.0"

  cluster_settings = var.ecs_cluster_settings
  cluster_name     = "${var.environment}-${var.unique_identifier}"
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
      essential = true
      image     = var.fargate_docker_image
      # Host port not required for Fargate tasks
      port_mappings = [
        {
          name          = "http"
          containerPort = var.fargate_container_port
          protocol      = "tcp"
          appProtocol   = "http"
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

  runtime_platform = {
    cpu_architecture        = var.fargate_cpu_architecture
    operating_system_family = "LINUX"
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

  # Wait for the database to be fully initialised before deploying app
  depends_on = [module.db]
}

# Used by ad-hoc ECS Fargate task to seed the database during E2E tests. Creating here to make it easier to pass config from Terraform / AWS
resource "aws_ecs_task_definition" "e2e_db_seeding" {
  family                   = "${var.service_name}-e2e-db-seeding"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.e2e_db_seed_task_cpu
  memory                   = var.e2e_db_seed_task_memory
  execution_role_arn       = module.ecs_service.task_exec_iam_role_arn

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = var.e2e_db_seed_cpu_architecture
  }

  container_definitions = jsonencode([
    {
      name      = "db-seeding"
      image     = var.e2e_db_seed_image
      essential = true
      environment = [
        {
          name : "RDS_USERNAME",
          value : var.db_master_username
        },
        {
          "name" : "RDS_ENDPOINT",
          "value" : module.db.cluster_endpoint
        },
        {
          "name" : "DB_NAME",
          "value" : var.db_database_name
        },
      ],

      secrets = [
        {
          "name" : "PGPASSWORD",
          "valueFrom" : aws_secretsmanager_secret_version.db_master_password.arn
        }
      ],

      logConfiguration = {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-group" : aws_cloudwatch_log_group.e2e_db_seeding.name,
          "awslogs-region" : var.region,
          "awslogs-stream-prefix" : "db-seeding"
        }
      },
    },
  ])
}