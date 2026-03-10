locals {
  secrets = {
    db_password          = var.db_password
    auth_api_secret      = var.auth_api_secret
    google_client_id     = var.google_client_id
    google_client_secret = var.google_client_secret
    facebook_app_id      = var.facebook_app_id
    facebook_app_secret  = var.facebook_app_secret
    sentry_dsn           = var.sentry_dsn
    stripe_secret_key    = var.stripe_secret_key
  }
}

resource "aws_secretsmanager_secret" "app" {
  for_each                = local.secrets
  name                    = "${var.project_name}/${var.environment}/${each.key}"
  recovery_window_in_days = 7

  tags = { Name = "${var.project_name}/${var.environment}/${each.key}" }
}

resource "aws_secretsmanager_secret_version" "app" {
  for_each      = local.secrets
  secret_id     = aws_secretsmanager_secret.app[each.key].id
  secret_string = each.value
}
