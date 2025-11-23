# Ham-Radio Cloud Architecture

## System Overview

Ham-Radio Cloud is a modern, cloud-native platform built as a monorepo with separate backend (Go) and frontend (SvelteKit) services.

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                            │
│                     SvelteKit + Tailwind                    │
│                    (PWA, Mobile-First)                      │
└────────────────────┬────────────────────────────────────────┘
                     │ REST/WebSocket
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend API                            │
│                     Go + Fiber                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │   Auth   │  │   QSO    │  │   Prop   │  │   SDR    │  │
│  │ Services │  │ Services │  │ Services │  │ Services │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
  ┌──────────┐  ┌────────┐  ┌─────────┐
  │PostgreSQL│  │ Redis  │  │ External│
  │+Timescale│  │ Cache  │  │  APIs   │
  └──────────┘  └────────┘  └─────────┘
```

## Directory Structure

```
/
├── backend/                 # Go backend service
│   ├── cmd/
│   │   └── api/            # API entry point
│   ├── internal/           # Private application code
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection & migrations
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Data models
│   │   └── services/       # Business logic
│   └── pkg/                # Public libraries
│
├── frontend/               # SvelteKit frontend
│   ├── src/
│   │   ├── lib/           # Shared components & utilities
│   │   │   ├── components/
│   │   │   ├── stores/
│   │   │   ├── utils/
│   │   │   └── types/
│   │   └── routes/        # SvelteKit routes
│   └── static/            # Static assets
│
├── shared/                # Shared types/utilities
│   └── types/            # Shared TypeScript types
│
├── infra/                # Infrastructure as code
│   ├── docker-compose.yml
│   └── init-scripts/
│
├── docs/                 # Documentation
│
└── scripts/              # Development scripts
```

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Fiber v2 (high-performance HTTP framework)
- **Database**: PostgreSQL 15 + TimescaleDB (time-series extension)
- **Cache**: Redis 7
- **Migrations**: golang-migrate

### Frontend
- **Framework**: SvelteKit (Svelte 5)
- **Styling**: Tailwind CSS
- **Build Tool**: Vite
- **PWA**: Service Worker + Manifest
- **TypeScript**: Strict mode enabled

### DevOps
- **Containerization**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Hosting** (Future): Fly.io or similar PaaS

## Data Models

### Core Entities

#### User
- Identity & authentication
- Callsign & QRZ verification
- Subscription tier (free/operator/contester)
- QSO limits & usage tracking

#### QSO (Contact Log)
- ADIF-compliant schema
- Time-series optimized (TimescaleDB hypertable)
- Full propagation & location data
- LoTW/eQSL sync status

#### Propagation Data (Future)
- Solar flux, Kp index
- VOACAP predictions
- Band condition forecasts

## API Design

### REST Endpoints

```
GET    /health                    # Health check
GET    /api/v1                    # API info

# Authentication
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/me

# QSOs
GET    /api/v1/qsos               # List QSOs (filtered)
POST   /api/v1/qsos               # Create QSO
GET    /api/v1/qsos/:id           # Get QSO
PUT    /api/v1/qsos/:id           # Update QSO
DELETE /api/v1/qsos/:id           # Delete QSO
POST   /api/v1/qsos/import        # ADIF import
GET    /api/v1/qsos/export        # ADIF export

# Propagation
GET    /api/v1/propagation        # Current conditions
GET    /api/v1/propagation/forecast

# SDR (Future)
GET    /api/v1/sdr/streams        # Public SDR streams
```

### WebSocket (Future)
- `/ws/spots` - Real-time DX cluster
- `/ws/propagation` - Live propagation updates

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    callsign VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    tier VARCHAR(20) DEFAULT 'free',
    qso_count INTEGER DEFAULT 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### QSOs Table (Hypertable)
```sql
CREATE TABLE qsos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    callsign VARCHAR(20) NOT NULL,
    time_on TIMESTAMP NOT NULL,
    band VARCHAR(10),
    mode VARCHAR(20),
    -- ... ADIF fields
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable('qsos', 'time_on');
```

## Authentication Flow

1. User clicks "Sign in with QRZ"
2. OAuth redirect to QRZ.com
3. QRZ validates & returns callsign
4. Backend creates/updates user
5. JWT token issued
6. Frontend stores token in localStorage
7. Subsequent requests include JWT in Authorization header

## Deployment Architecture (Future)

### MVP Deployment
- **Backend**: Fly.io (multi-region)
- **Frontend**: Vercel or Netlify
- **Database**: Managed PostgreSQL (Fly.io or Supabase)
- **Cache**: Redis on Fly.io

### Production Scaling
- CDN for static assets
- Read replicas for database
- Horizontal scaling with load balancer
- Background job queue (Redis + Worker pods)

## Security Considerations

### Authentication
- JWT with short expiration (15min)
- Refresh tokens (7 days)
- MFA for operator+ tiers

### Authorization
- Role-based access control (RBAC)
- API rate limiting (per tier)
- QSO limits enforced at DB & API layers

### Data Protection
- TLS/HTTPS everywhere
- Encrypted database backups
- PII minimization (callsigns are public)
- GDPR compliance (EU users)

## Observability

### Logging
- Structured JSON logs
- Log levels: DEBUG, INFO, WARN, ERROR
- Correlation IDs for request tracking

### Metrics (Future)
- Prometheus + Grafana
- API latency, error rates
- QSO throughput, DB connection pool

### Monitoring (Future)
- Uptime checks (external)
- Alert on critical failures
- Database performance monitoring

## Development Workflow

1. **Local Setup**: `./scripts/setup.sh`
2. **Start Services**: `./scripts/dev.sh`
3. **Run Tests**: `go test ./...` (backend), `npm test` (frontend)
4. **Lint**: `npm run lint` (frontend), `golangci-lint run` (backend)
5. **Commit**: Git hooks run linters
6. **Push**: CI runs tests + builds
7. **Deploy**: Manual deploy to staging → production

## Phase 1 Deliverables

- ✅ Monorepo structure
- ✅ Backend skeleton (Go + Fiber)
- ✅ Frontend skeleton (SvelteKit + Tailwind)
- ✅ Database schema (users + QSOs)
- ✅ Docker Compose (local dev)
- ✅ CI/CD (GitHub Actions)
- 🔲 OAuth (QRZ.com)
- 🔲 QSO CRUD APIs
- 🔲 ADIF import/export
- 🔲 Basic propagation display

---

**Next**: See [CONTRIBUTING.md](./CONTRIBUTING.md) for development guidelines.
