resource "aws_secretsmanager_secret" "db_master_password" {
  name        = "${var.environment}/${var.service_name}/db_master_password"
  description = "Master password for the ${module.db.cluster_id} Aurora cluster"

  # Allow the E2E tests to repeatedly create/delete the same key name
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "db_master_password" {
  secret_id     = aws_secretsmanager_secret.db_master_password.id
  secret_string = random_password.aurora_master_password.result
}