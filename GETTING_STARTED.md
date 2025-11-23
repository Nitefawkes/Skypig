# Getting Started with Ham-Radio Cloud

Welcome to Ham-Radio Cloud! This guide will get you up and running in minutes.

## Prerequisites

Make sure you have these installed:

- **Docker** & Docker Compose - [Install Docker](https://docs.docker.com/get-docker/)
- **Go** 1.21+ - [Install Go](https://go.dev/doc/install)
- **Node.js** 20+ - [Install Node](https://nodejs.org/)
- **Git** - [Install Git](https://git-scm.com/)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/Nitefawkes/Skypig.git
cd Skypig
```

### 2. Run Setup Script

This will:
- Copy environment variable templates
- Start PostgreSQL & Redis via Docker
- Install dependencies

```bash
./scripts/setup.sh
```

### 3. Start Development Servers

```bash
./scripts/dev.sh
```

This starts both backend and frontend in development mode.

### 4. Open Your Browser

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/health
- **Database Admin**: http://localhost:8081

## Project Structure

```
Skypig/
├── backend/         # Go API (Fiber framework)
├── frontend/        # SvelteKit app (Tailwind CSS)
├── shared/          # Shared TypeScript types
├── infra/           # Docker Compose & infrastructure
├── docs/            # Documentation
└── scripts/         # Helper scripts
```

## What's Included

### Backend (Go + Fiber)
- RESTful API with health checks
- PostgreSQL database with TimescaleDB
- User and QSO data models
- Database migrations
- Environment-based configuration

### Frontend (SvelteKit + Tailwind)
- Modern UI with Tailwind CSS
- Progressive Web App (PWA) ready
- API integration with backend
- Mobile-first responsive design

### Infrastructure
- Docker Compose for local development
- PostgreSQL 15 + TimescaleDB extension
- Redis for caching
- Adminer for database management

## Running Individually

### Backend Only

```bash
cd backend
go run cmd/api/main.go
```

Runs on http://localhost:8080

### Frontend Only

```bash
cd frontend
npm run dev
```

Runs on http://localhost:5173

### Database Only

```bash
cd infra
docker-compose up -d postgres
```

PostgreSQL on localhost:5432 (user: `postgres`, password: `postgres`, db: `hamradio_cloud`)

## Development Workflow

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes

3. Test your changes:
   ```bash
   # Backend tests
   cd backend && go test ./...

   # Frontend checks
   cd frontend && npm run check && npm run lint
   ```

4. Commit and push:
   ```bash
   git add .
   git commit -m "feat: add my feature"
   git push origin feature/my-feature
   ```

### Testing the API

```bash
# Health check
curl http://localhost:8080/health

# API info
curl http://localhost:8080/api/v1
```

## Common Commands

```bash
# Format backend code
cd backend && go fmt ./...

# Format frontend code
cd frontend && npm run format

# View Docker logs
docker-compose -f infra/docker-compose.yml logs -f

# Access database CLI
docker exec -it hamradio_postgres psql -U postgres -d hamradio_cloud

# Stop all services
docker-compose -f infra/docker-compose.yml down
```

## Environment Variables

Key configuration files (copy from `.env.example` files):

- `/.env` - Root environment variables
- `/backend/.env` - Backend configuration
- `/frontend/.env` - Frontend configuration

### Important Backend Variables

```bash
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hamradio_cloud
```

### Important Frontend Variables

```bash
PUBLIC_API_URL=http://localhost:8080
```

## Troubleshooting

### Docker Services Won't Start

```bash
cd infra
docker-compose down -v
docker-compose up -d
```

### Backend Won't Connect to Database

1. Check Docker is running: `docker ps`
2. Verify PostgreSQL is healthy: `docker-compose -f infra/docker-compose.yml logs postgres`
3. Check environment variables in `backend/.env`

### Frontend Build Errors

```bash
cd frontend
rm -rf node_modules .svelte-kit
npm install
npm run dev
```

## Next Steps

Now that you're set up:

1. **Read the docs**:
   - [ARCHITECTURE.md](docs/ARCHITECTURE.md) - System design
   - [CONTRIBUTING.md](docs/CONTRIBUTING.md) - Development guidelines
   - [DEVELOPMENT.md](docs/DEVELOPMENT.md) - Developer guide

2. **Explore the code**:
   - Check out `backend/cmd/api/main.go` for the API entry point
   - Look at `frontend/src/routes/+page.svelte` for the homepage

3. **Start building**:
   - Phase 2 tasks are in the [Product Blueprint](README.md)
   - Next up: OAuth integration, QSO CRUD, ADIF import/export

## Resources

- [Product Blueprint](README.md) - Original project plan
- [Go Documentation](https://go.dev/doc/)
- [SvelteKit Docs](https://kit.svelte.dev/)
- [Fiber Framework](https://docs.gofiber.io/)
- [Tailwind CSS](https://tailwindcss.com/)

## Need Help?

- Check the [docs/](docs/) folder
- Open an issue on GitHub
- Review the troubleshooting section above

---

**Welcome aboard! 73s and happy coding!** 📻
