#!/bin/bash

# Exit on any error
set -e

# Check if .env.secret exists
if [ ! -f .env.secret ]; then
    echo "Error: .env.secret file not found"
    exit 1
fi

# Source the secret environment variables
source .env.secret

# Create the namespace if it doesn't exist
kubectl create namespace otel-demo --dry-run=client -o yaml | kubectl apply -f -

# Create the secret
kubectl create secret generic product-catalog-db-secret \
    --namespace otel-demo \
    --from-literal=database-url="$NEON_DB_URL" \
    --dry-run=client -o yaml | kubectl apply -f -

echo "Secrets applied successfully" 