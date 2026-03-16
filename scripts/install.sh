#!/bin/bash

# DaaS Backend Installation Script
# This script prepares the environment and starts the containers.

set -e

echo "🚀 Starting DaaS Backend Installation..."

# 1. Copy .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
    echo "⚠️  Please update .env with your actual credentials if necessary."
fi

# 4. Pull or build the containers
echo "🛠 Building and starting containers..."
docker-compose up -d --build

# 5. Run DB setup script inside the container (optional manually or here)
echo "📂 Setting up database..."
chmod +x scripts/setup-db.sh
./scripts/setup-db.sh

echo "✅ Installation complete! API is running on http://localhost:8084"
echo "🔍 Check logs with: docker-compose logs -f"
