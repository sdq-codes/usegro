#!/bin/sh
set -e

envsubst < /app/config/config.prod.yaml > /app/config/config.prod.yaml.tmp
mv /app/config/config.prod.yaml.tmp /app/config/config.prod.yaml

exec "$@"
