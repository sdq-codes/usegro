variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "eu-west-1"
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "usegro"
}

variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "production"
}

# ─── Network ────────────────────────────────────────────────────────────────

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones to use"
  type        = list(string)
  default     = ["eu-west-1a", "eu-west-1b"]
}

# ─── RDS ────────────────────────────────────────────────────────────────────

variable "db_name" {
  description = "PostgreSQL database name"
  type        = string
  default     = "usegro"
}

variable "db_username" {
  description = "PostgreSQL master username"
  type        = string
  default     = "usegro"
}

variable "db_password" {
  description = "PostgreSQL master password"
  type        = string
  sensitive   = true
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

# ─── ElastiCache ────────────────────────────────────────────────────────────

variable "redis_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t3.micro"
}

# ─── Secrets ────────────────────────────────────────────────────────────────

variable "auth_api_secret" {
  description = "JWT / auth API secret"
  type        = string
  sensitive   = true
}

variable "google_client_id" {
  description = "Google OAuth client ID"
  type        = string
  sensitive   = true
  default     = ""
}

variable "google_client_secret" {
  description = "Google OAuth client secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "facebook_app_id" {
  description = "Facebook app ID"
  type        = string
  sensitive   = true
  default     = ""
}

variable "facebook_app_secret" {
  description = "Facebook app secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "sentry_dsn" {
  description = "Sentry DSN"
  type        = string
  sensitive   = true
  default     = ""
}

variable "ses_from_email" {
  description = "SES sender email address"
  type        = string
  default     = ""
}

variable "stripe_secret_key" {
  description = "Stripe secret key"
  type        = string
  sensitive   = true
  default     = ""
}


# ─── ECS ────────────────────────────────────────────────────────────────────

variable "base_desired_count" {
  description = "Desired task count for base service"
  type        = number
  default     = 1
}

variable "crm_desired_count" {
  description = "Desired task count for crm service"
  type        = number
  default     = 1
}

variable "catalog_desired_count" {
  description = "Desired task count for catalog service"
  type        = number
  default     = 1
}

variable "billing_desired_count" {
  description = "Desired task count for billing service"
  type        = number
  default     = 1
}


# ─── Domain ─────────────────────────────────────────────────────────────────

variable "domain_name" {
  description = "Primary domain name (e.g. api.usegro.com). Leave empty to skip Route53/ACM."
  type        = string
  default     = ""
}

variable "certificate_arn" {
  description = "ACM certificate ARN for HTTPS. Leave empty to use HTTP only."
  type        = string
  default     = ""
}
