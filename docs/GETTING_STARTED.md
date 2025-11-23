# Getting Started with Ham Radio Cloud

This guide will help you set up and run the Ham Radio Cloud development environment.

## Prerequisites

- **Go** 1.23+ (for backend)
- **Node.js** 20+ (for frontend)
- **Docker** & **Docker Compose** (for local database and services)
- **Make** (optional, for convenience commands)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/nitefawkes/ham-radio-cloud.git
cd ham-radio-cloud
```

### 2. Environment Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` and add your configuration:
- QRZ.com OAuth credentials (get from https://www.qrz.com/docs/logbook-api)
- JWT secret (use a strong random string)
- Database credentials (default works for local development)

### 3. Start the Development Environment

#### Option A: Using Make (Recommended)

```bash
# Install all dependencies
make install

# Start all services (Docker Compose)
make dev
```

#### Option B: Using Docker Compose Directly

```bash
# Install dependencies
cd backend && go mod download && cd ..
cd frontend && npm install && cd ..

# Start services
docker-compose -f deployments/docker/docker-compose.yml up
```

#### Option C: Manual Setup (Without Docker)

**Terminal 1 - Database:**
```bash
# Install PostgreSQL with TimescaleDB locally or use a cloud provider
# Run the init.sql script to set up the schema
psql -U postgres -d hamradio -f deployments/docker/init.sql
```

**Terminal 2 - Backend:**
```bash
cd backend
air  # or: go run cmd/api/main.go
```

**Terminal 3 - Frontend:**
```bash
cd frontend
npm run dev
```

### 4. Access the Application

- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **API Health Check:** http://localhost:8080/api/v1/health

## Development Workflow

### Backend Development

```bash
# Run backend with hot reload
make dev-backend

# Run tests
cd backend && go test -v ./...

# Format code
cd backend && gofmt -s -w .

# Lint
cd backend && go vet ./...
```

### Frontend Development

```bash
# Run frontend dev server
make dev-frontend

# Type check
cd frontend && npm run check

# Lint
cd frontend && npm run lint

# Format
cd frontend && npm run format
```

### Database Operations

```bash
# View database logs
make docker-logs

# Reset database (WARNING: destroys all data)
make db-reset

# Run migrations
make db-migrate
```

## Project Structure

```
ham-radio-cloud/
├── backend/              # Go + Fiber backend
│   ├── cmd/api/          # Main application entry point
│   ├── internal/         # Internal packages
│   │   ├── handlers/     # HTTP handlers
│   │   ├── services/     # Business logic
│   │   ├── repositories/ # Data access layer
│   │   ├── models/       # Data models
│   │   ├── middleware/   # HTTP middleware
│   │   └── config/       # Configuration
│   └── pkg/              # Public packages
├── frontend/             # SvelteKit frontend
│   ├── src/
│   │   ├── routes/       # Pages and API routes
│   │   ├── lib/          # Shared components, stores, utils
│   │   └── app.css       # Global styles
│   └── static/           # Static assets
├── deployments/          # Deployment configurations
│   └── docker/           # Docker & Docker Compose files
├── docs/                 # Documentation
├── scripts/              # Utility scripts
└── shared/               # Shared code between frontend/backend
```

## Next Steps

1. Read the [Architecture Overview](./ARCHITECTURE.md)
2. Review the [API Documentation](./API.md)
3. Check the [Product Blueprint](../README.md)
4. Join the development workflow in [CONTRIBUTING.md](./CONTRIBUTING.md)

## Troubleshooting

### Port Already in Use

If ports 8080 or 5173 are already in use, you can change them in:
- Backend: `.env` file (`PORT` variable)
- Frontend: `frontend/vite.config.ts` (server.port)

### Database Connection Issues

Make sure PostgreSQL is running:
```bash
docker-compose -f deployments/docker/docker-compose.yml ps
```

Check logs:
```bash
docker-compose -f deployments/docker/docker-compose.yml logs postgres
```

### Build Errors

Clean and reinstall:
```bash
make clean
make install
```

## Getting Help

- Check the [FAQ](./FAQ.md)
- Review existing [GitHub Issues](https://github.com/nitefawkes/ham-radio-cloud/issues)
- Join our community discussions

---

73 de W1AW 📻
