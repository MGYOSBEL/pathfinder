#!/usr/bin/env bash
set -euo pipefail

NETWORK_NAME="pathfinder_network_dev"

echo "🔍  Checking if network '$NETWORK_NAME' exists..."

if docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
  echo "✅  Network '$NETWORK_NAME' already exists. Skipping creation."
else
  echo "⚡  Network '$NETWORK_NAME' not found. Creating..."
  docker network create "$NETWORK_NAME"
  echo "🎉  Network '$NETWORK_NAME' created successfully."
fi
