#!/usr/bin/env bash
# Generates config.prod.yaml for base and crm services by filling in
# placeholder values from terraform outputs and terraform.tfvars.
#
# Usage: ./infrastructure/scripts/generate-configs.sh
# Run this from the repo root AFTER terraform apply.

set -euo pipefail

TERRAFORM_DIR="infrastructure/terraform"
BASE_CONFIG="services/base/config/config.prod.yaml"
CRM_CONFIG="services/crm/config/config.prod.yaml"

echo "→ Reading terraform outputs..."
cd "$TERRAFORM_DIR"

RDS_HOST=$(terraform output -raw rds_endpoint | cut -d: -f1)
REDIS_HOST=$(terraform output -raw redis_primary_endpoint)
ALB_DNS=$(terraform output -raw alb_dns_name)

# Read static values from terraform.tfvars
AWS_REGION=$(grep 'aws_region' terraform.tfvars | awk -F'"' '{print $2}')
DB_NAME=$(grep 'db_name' terraform.tfvars | awk -F'"' '{print $2}')
DB_USERNAME=$(grep 'db_username' terraform.tfvars | awk -F'"' '{print $2}')
DOMAIN_NAME=$(grep 'domain_name' terraform.tfvars | awk -F'"' '{print $2}')

cd - > /dev/null

# Use domain if set, otherwise fall back to ALB DNS
if [ -n "$DOMAIN_NAME" ]; then
  ALB_URL="https://${DOMAIN_NAME}"
  FRONTEND_URL="https://${DOMAIN_NAME}"
else
  ALB_URL="http://${ALB_DNS}"
  FRONTEND_URL="http://${ALB_DNS}"
fi

echo "  RDS host    : $RDS_HOST"
echo "  Redis host  : $REDIS_HOST"
echo "  ALB URL     : $ALB_URL"
echo "  AWS region  : $AWS_REGION"
echo "  DB name     : $DB_NAME"
echo ""

fill_placeholders() {
  local file="$1"
  sed -i.bak \
    -e "s|__RDS_HOST__|${RDS_HOST}|g" \
    -e "s|__REDIS_HOST__|${REDIS_HOST}|g" \
    -e "s|__ALB_URL__|${ALB_URL}|g" \
    -e "s|__FRONTEND_URL__|${FRONTEND_URL}|g" \
    -e "s|__AWS_REGION__|${AWS_REGION}|g" \
    -e "s|__DB_NAME__|${DB_NAME}|g" \
    -e "s|__DB_USERNAME__|${DB_USERNAME}|g" \
    "$file"
  rm -f "${file}.bak"
  echo "✓ $file"
}

fill_placeholders "$BASE_CONFIG"
fill_placeholders "$CRM_CONFIG"

echo ""
echo "Done. Build your Docker images now:"
echo "  docker build -f infrastructure/docker/base.Dockerfile -t \$BASE_ECR:production ."
echo "  docker build -f infrastructure/docker/crm.Dockerfile  -t \$CRM_ECR:production ."
