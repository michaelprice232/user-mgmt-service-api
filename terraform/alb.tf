module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 9.0"

  name                       = substr(local.name, 0, 32) # API limit of 32
  load_balancer_type         = "application"
  vpc_id                     = module.vpc.vpc_id
  subnets                    = module.vpc.public_subnets
  enable_deletion_protection = false

  security_group_name            = "${local.name}-lb"
  security_group_use_name_prefix = false
  security_group_description     = "Assigned to ${local.name} load balancer"
  security_group_ingress_rules = {
    all_http = {
      from_port   = 80
      to_port     = 80
      ip_protocol = "tcp"
      cidr_ipv4   = "0.0.0.0/0"
    }
  }
  security_group_egress_rules = {
    all = {
      ip_protocol = "-1"
      cidr_ipv4   = module.vpc.vpc_cidr_block
    }
  }
  security_group_tags = {
    Name : "${local.name}-lb"
  }

  listeners = {
    http = {
      port     = 80
      protocol = "HTTP"

      forward = {
        target_group_key = "ecs"
      }
    }
  }

  target_groups = {
    ecs = {
      backend_protocol                  = "HTTP"
      port                              = 8080
      target_type                       = "ip"
      deregistration_delay              = 5
      load_balancing_cross_zone_enabled = true

      health_check = {
        enabled             = true
        healthy_threshold   = 2
        interval            = 30
        matcher             = "200"
        path                = "/health"
        port                = "traffic-port"
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 2
      }

      # There's nothing to attach here in this definition. Instead,
      # ECS will attach the IPs of the tasks to this target group
      create_attachment = false
    }
  }
}