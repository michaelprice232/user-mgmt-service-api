# Stores logs from the DB seeding container used during E2E tests in AWS
# Other log groups for the app are created the ecs_service Terraform module
resource "aws_cloudwatch_log_group" "e2e_db_seeding" {
  name              = "/aws/ecs/${var.service_name}/e2e-db-seeding"
  retention_in_days = var.e2e_db_seed_log_retention
}