# Contributing to Ham-Radio Cloud

Thank you for your interest in contributing to Ham-Radio Cloud! This document provides guidelines for development.

## Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Node.js** 20 or higher
- **Docker** & Docker Compose
- **Git**

### Initial Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Nitefawkes/Skypig.git
   cd Skypig
   ```

2. Run the setup script:
   ```bash
   ./scripts/setup.sh
   ```

3. Start development servers:
   ```bash
   ./scripts/dev.sh
   ```

4. Access the application:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Database Admin: http://localhost:8081

## Development Workflow

### Branch Strategy

- `main` - Production-ready code
- `develop` - Integration branch for features
- `feature/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation updates

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the code style guidelines below

3. Test your changes:
   ```bash
   # Backend
   cd backend
   go test ./...

   # Frontend
   cd frontend
   npm run check
   npm run lint
   ```

4. Commit with a clear message:
   ```bash
   git commit -m "Add feature: brief description"
   ```

5. Push and create a pull request:
   ```bash
   git push origin feature/your-feature-name
   ```

## Code Style Guidelines

### Go (Backend)

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new functionality
- Keep functions small and focused
- Use meaningful variable names

```go
// Good
func CreateQSO(ctx context.Context, qso *models.QSO) error {
    // Implementation
}

// Bad
func cq(c context.Context, q *models.QSO) error {
    // Implementation
}
```

### TypeScript/Svelte (Frontend)

- Use TypeScript strict mode
- Follow the Prettier configuration
- Use Tailwind utility classes (avoid custom CSS)
- Keep components small and reusable
- Use `$state` and `$derived` for reactivity (Svelte 5)

```svelte
<!-- Good -->
<script lang="ts">
	import type { QSO } from '$types';

	let qsos = $state<QSO[]>([]);
</script>

<div class="card">
	{#each qsos as qso}
		<QSOCard {qso} />
	{/each}
</div>
```

### Database Migrations

- Always create both `.up.sql` and `.down.sql` files
- Use sequential numbering: `001_`, `002_`, etc.
- Test rollbacks before committing
- Document complex migrations

```sql
-- 003_add_user_preferences.up.sql
ALTER TABLE users ADD COLUMN preferences JSONB DEFAULT '{}';
CREATE INDEX idx_users_preferences ON users USING GIN (preferences);
```

## Testing

### Backend Tests

```bash
cd backend
go test ./...                    # Run all tests
go test -v ./internal/handlers  # Verbose specific package
go test -race ./...              # Race detection
go test -cover ./...             # Coverage report
```

### Frontend Tests

```bash
cd frontend
npm run check       # Type checking
npm run lint        # Linting
npm test            # Unit tests (when implemented)
```

## Project Structure Conventions

### Backend (`/backend`)

```
internal/
  handlers/     - HTTP request handlers
  services/     - Business logic (no HTTP dependencies)
  models/       - Data structures
  middleware/   - HTTP middleware
  database/     - DB connection & queries
  config/       - Configuration management
```

**Rules:**
- `internal/handlers` should be thin, delegating to `services`
- `services` should not import `handlers` or `fiber`
- All database queries should be in `database/` or `services/`

### Frontend (`/frontend`)

```
src/
  lib/
    components/  - Reusable UI components
    stores/      - Svelte stores (global state)
    utils/       - Helper functions
    types/       - TypeScript types
  routes/        - SvelteKit routes (pages)
```

**Rules:**
- Components should be in `lib/components`
- Pages should be in `routes`
- Avoid prop drilling - use stores for deep state
- Keep API calls in `lib/api` (to be created)

## API Design Principles

### RESTful Endpoints

- Use plural nouns: `/api/v1/qsos`, not `/api/v1/qso`
- Use HTTP methods correctly:
  - `GET` - Retrieve data
  - `POST` - Create new resource
  - `PUT` - Update entire resource
  - `PATCH` - Update partial resource
  - `DELETE` - Remove resource

### Response Format

```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 1234
  }
}
```

### Error Format

```json
{
  "error": {
    "code": "INVALID_CALLSIGN",
    "message": "Callsign must be 3-10 alphanumeric characters",
    "details": { ... }
  }
}
```

## Database Guidelines

### Naming Conventions

- Tables: plural, lowercase, snake_case (`users`, `qso_logs`)
- Columns: singular, lowercase, snake_case (`callsign`, `time_on`)
- Indexes: `idx_tablename_columnname`
- Foreign keys: `fk_tablename_columnname`

### Performance Considerations

- Index frequently queried columns
- Use `EXPLAIN ANALYZE` for slow queries
- Leverage TimescaleDB for time-series data (QSOs)
- Avoid N+1 queries (use JOINs or batch loading)

## Commit Message Guidelines

Format: `<type>: <subject>`

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Formatting, no code change
- `refactor` - Code restructuring
- `test` - Adding tests
- `chore` - Maintenance tasks

Examples:
```
feat: add ADIF import endpoint
fix: correct QSO date parsing in UTC
docs: update API documentation for /qsos endpoint
refactor: extract QSO validation into service layer
```

## Pull Request Process

1. Update relevant documentation
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md (if applicable)
5. Request review from maintainers
6. Address review feedback
7. Squash commits if requested
8. Merge after approval

## Development Tips

### Hot Reload

Both backend and frontend support hot reload:

- **Backend**: Use `air` for live reload (optional)
  ```bash
  go install github.com/cosmtrek/air@latest
  cd backend && air
  ```

- **Frontend**: Vite provides instant HMR
  ```bash
  cd frontend && npm run dev
  ```

### Database Inspection

Use Adminer at http://localhost:8081 or:

```bash
docker exec -it hamradio_postgres psql -U postgres -d hamradio_cloud
```

### Debugging

- **Backend**: Use VS Code Go debugger or `dlv`
- **Frontend**: Chrome/Firefox DevTools
- **Database**: `EXPLAIN ANALYZE` for query plans

## Common Issues

### Docker services won't start

```bash
cd infra
docker-compose down -v  # Remove volumes
docker-compose up -d
```

### Go dependencies out of sync

```bash
cd backend
go mod tidy
go mod download
```

### Frontend build errors

```bash
cd frontend
rm -rf node_modules .svelte-kit
npm install
```

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [SvelteKit Documentation](https://kit.svelte.dev/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [TimescaleDB Docs](https://docs.timescale.com/)

## Questions?

- Open an issue on GitHub
- Check existing documentation in `/docs`
- Review the [ARCHITECTURE.md](./ARCHITECTURE.md) for system design

---

**Thank you for contributing to Ham-Radio Cloud!** 73s!
