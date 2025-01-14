output "service_endpoint" {
  value = format("http://%s", module.alb.dns_name)
}

output "service_name" {
  value = module.ecs_service.name
}

output "ecs_cluster_name" {
  value = module.ecs_cluster.cluster_name
}

output "private_subnets" {
  value = module.vpc.private_subnets
}

output "public_subnets" {
  value = module.vpc.public_subnets
}

output "fargate_task_security_group_id" {
  value = module.ecs_service.security_group_id
}

output "db_seeding_task_definition_target" {
  value = format("%s:%s", aws_ecs_task_definition.e2e_db_seeding.family, aws_ecs_task_definition.e2e_db_seeding.revision)
}

output "db_seeding_target_subnet" {
  value = module.vpc.private_subnets[0]
}