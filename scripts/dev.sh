#!/bin/bash

# Ham-Radio Cloud Development Script
# Starts both backend and frontend in development mode

set -e

echo "🚀 Starting Ham-Radio Cloud in development mode..."

# Check if Docker services are running
if ! docker ps | grep -q hamradio_postgres; then
    echo "📦 Starting Docker services..."
    cd infra
    docker-compose up -d
    cd ..
    echo "⏳ Waiting for services to be ready..."
    sleep 3
fi

# Start backend in background
echo "🔧 Starting backend API..."
cd backend
go run cmd/api/main.go &
BACKEND_PID=$!
cd ..

# Give backend time to start
sleep 2

# Start frontend
echo "🎨 Starting frontend..."
cd frontend
npm run dev

# Cleanup on exit
trap "echo '🛑 Shutting down...'; kill $BACKEND_PID 2>/dev/null; exit" INT TERM
