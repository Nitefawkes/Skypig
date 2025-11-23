#!/bin/bash
set -e

echo "🚀 Setting up Ham-Radio Cloud development environment..."

# Check for required tools
command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || command -v docker compose >/dev/null 2>&1 || { echo "❌ Docker Compose is required but not installed. Aborting." >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Aborting." >&2; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is required but not installed. Aborting." >&2; exit 1; }

echo "✅ All required tools are installed"

# Copy environment files
if [ ! -f .env ]; then
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
fi

if [ ! -f backend/.env ]; then
    echo "📝 Creating backend/.env from backend/.env.example..."
    cp backend/.env.example backend/.env
fi

if [ ! -f frontend/.env ]; then
    echo "📝 Creating frontend/.env from frontend/.env.example..."
    cp frontend/.env.example frontend/.env
fi

# Start Docker services
echo "🐳 Starting Docker services..."
cd infra
docker-compose up -d
cd ..

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Install backend dependencies
echo "📦 Installing Go dependencies..."
cd backend
go mod download
cd ..

# Install frontend dependencies
echo "📦 Installing Node.js dependencies..."
cd frontend
npm install
cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Start the backend:  cd backend && go run cmd/api/main.go"
echo "  2. Start the frontend: cd frontend && npm run dev"
echo "  3. Open http://localhost:5173 in your browser"
echo ""
echo "Database is running on localhost:5432"
echo "Adminer (DB UI) is available at http://localhost:8081"
echo ""
