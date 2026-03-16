#!/bin/bash

# DaaS Backend Installation Script
# This script prepares the environment and starts the containers.

set -e

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# Project root is one level up
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "🚀 Starting DaaS Backend Installation from $PROJECT_ROOT..."

# 1. Copy .env if it doesn't exist (non-fatal if permission denied)
if [ ! -f .env ]; then
    echo "📝 Attempting to create .env from .env.example..."
    cp .env.example .env 2>/dev/null || echo "⚠️ Could not copy .env (possibly already exists or permission denied). Skipping..."
fi

# 4. Pull or build the containers
echo "🛠 Building and starting containers..."
docker compose up -d --build

# 5. Run DB setup script inside the container (optional manually or here)
echo "📂 Setting up database..."
chmod +x scripts/setup-db.sh
./scripts/setup-db.sh

echo "✅ Installation complete! API is running on http://localhost:8084"
echo "🔍 Check logs with: docker compose logs -f"
