# Ham-Radio Cloud - Development Guide

Quick start guide for developers working on Ham-Radio Cloud.

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/Nitefawkes/Skypig.git
cd Skypig

# 2. Run setup (creates .env files, starts Docker)
./scripts/setup.sh

# 3. Start development servers
./scripts/dev.sh
```

Access the app at http://localhost:5173

## Project Structure

```
Skypig/
├── backend/         # Go + Fiber API
├── frontend/        # SvelteKit + Tailwind
├── shared/          # Shared types
├── infra/           # Docker Compose
├── docs/            # Documentation
└── scripts/         # Helper scripts
```

## Running Services Individually

### Backend Only

```bash
cd backend
go run cmd/api/main.go
```

Backend runs on http://localhost:8080

### Frontend Only

```bash
cd frontend
npm run dev
```

Frontend runs on http://localhost:5173 (proxies `/api` to backend)

### Database Only

```bash
cd infra
docker-compose up postgres
```

PostgreSQL runs on localhost:5432

## Testing

```bash
# Backend tests
cd backend && go test ./...

# Frontend checks
cd frontend && npm run check && npm run lint

# Run all tests
./scripts/test.sh  # (to be created)
```

## Database Migrations

Migrations are in `backend/internal/database/migrations/`

```bash
# Apply migrations (manual for now)
# TODO: Add migration tool
```

## Environment Variables

Copy `.env.example` files and customize:

```bash
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

## Common Tasks

### Add a new API endpoint

1. Create handler in `backend/internal/handlers/`
2. Add route in `backend/cmd/api/main.go`
3. Test with curl or Postman
4. Add frontend API call in `frontend/src/lib/api/`

### Add a new page

1. Create route in `frontend/src/routes/`
2. Use existing components from `frontend/src/lib/components/`
3. Add navigation link if needed

### Create database migration

1. Add files in `backend/internal/database/migrations/`:
   - `00X_description.up.sql`
   - `00X_description.down.sql`
2. Test up and down migrations
3. Update models in `backend/internal/models/`

## Useful Commands

```bash
# Format Go code
cd backend && go fmt ./...

# Format frontend code
cd frontend && npm run format

# Build for production
cd backend && go build -o api cmd/api/main.go
cd frontend && npm run build

# View Docker logs
docker-compose -f infra/docker-compose.yml logs -f

# Access database
docker exec -it hamradio_postgres psql -U postgres -d hamradio_cloud
```

## Debugging

### Backend
- Use VS Code with Go extension
- Or use `dlv` debugger: `dlv debug cmd/api/main.go`

### Frontend
- Use browser DevTools
- Svelte DevTools extension for Chrome/Firefox

### Database
- Adminer UI: http://localhost:8081
- Or psql: `docker exec -it hamradio_postgres psql -U postgres`

## Phase 1 Checklist

### Infrastructure ✅
- [x] Monorepo structure
- [x] Go backend with Fiber
- [x] SvelteKit frontend with Tailwind
- [x] PostgreSQL + TimescaleDB
- [x] Docker Compose
- [x] CI/CD (GitHub Actions)

### Next Steps 🔲
- [ ] Implement OAuth (QRZ.com)
- [ ] QSO CRUD endpoints
- [ ] ADIF import/export
- [ ] Basic propagation display
- [ ] User authentication flow
- [ ] Dashboard UI

## Resources

- [Product Blueprint](../README.md) - Original project plan
- [Architecture](./ARCHITECTURE.md) - System design
- [Contributing](./CONTRIBUTING.md) - Development guidelines

## Need Help?

- Check the [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines
- Open an issue on GitHub
- Review [ARCHITECTURE.md](./ARCHITECTURE.md) for system design questions

---

Happy coding! 73s!
