# Ham Radio Cloud - Architecture Overview

## System Architecture

Ham Radio Cloud follows a modern, cloud-native architecture designed for scalability, maintainability, and performance.

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Browser    │  │  Mobile PWA  │  │  API Clients │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Frontend (SvelteKit)                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Routes    Components    Stores    Utils    Types        │  │
│  └──────────────────────────────────────────────────────────┘  │
│  • SSR/CSR Hybrid • PWA • Tailwind CSS • TypeScript          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼ REST API
┌─────────────────────────────────────────────────────────────────┐
│                     Backend (Go + Fiber)                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              API Layer (v1)                               │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐         │  │
│  │  │  Handlers  │→ │  Services  │→ │Repositories│         │  │
│  │  └────────────┘  └────────────┘  └────────────┘         │  │
│  └──────────────────────────────────────────────────────────┘  │
│  • Clean Architecture • Middleware • JWT Auth • CORS          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Database (PostgreSQL + TimescaleDB)           │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐              │
│  │   Users    │  │    QSOs    │  │Propagation │              │
│  └────────────┘  └────────────┘  └────────────┘              │
│  • Time-Series Hypertables • Indexes • Triggers               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     External Services                           │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐              │
│  │  QRZ OAuth │  │NOAA SWPC   │  │   LoTW     │              │
│  └────────────┘  └────────────┘  └────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Backend Architecture (Go)

### Clean Architecture Layers

```
cmd/api/
  └── main.go              # Application entry point

internal/
  ├── handlers/            # HTTP request handlers
  │   ├── auth.go
  │   ├── qso.go
  │   ├── propagation.go
  │   └── user.go
  │
  ├── services/            # Business logic layer
  │   ├── qso_service.go
  │   ├── lotw_service.go
  │   └── propagation_service.go
  │
  ├── repositories/        # Data access layer
  │   ├── qso_repository.go
  │   └── user_repository.go
  │
  ├── models/              # Data models
  │   ├── user.go
  │   ├── qso.go
  │   └── propagation.go
  │
  ├── middleware/          # HTTP middleware
  │   ├── auth.go
  │   ├── rate_limit.go
  │   └── logging.go
  │
  └── config/              # Configuration
      └── config.go

pkg/                       # Public packages
  ├── adif/                # ADIF parser/exporter
  ├── lotw/                # LoTW integration
  └── voacap/              # VOACAP integration
```

### Request Flow

```
HTTP Request
    ↓
[Middleware Stack]
    ↓ (Logging, CORS, Auth)
[Handler Layer]
    ↓ (Validation, HTTP concerns)
[Service Layer]
    ↓ (Business logic)
[Repository Layer]
    ↓ (Database queries)
[Database]
```

---

## Frontend Architecture (SvelteKit)

### Directory Structure

```
src/
  ├── routes/              # File-based routing
  │   ├── +layout.svelte   # Root layout
  │   ├── +page.svelte     # Home page
  │   ├── logbook/
  │   ├── propagation/
  │   └── sdr/
  │
  ├── lib/
  │   ├── components/      # Reusable UI components
  │   ├── stores/          # Svelte stores (state)
  │   ├── utils/           # Utilities
  │   │   └── api.ts       # API client
  │   └── types/           # TypeScript types
  │
  ├── service-worker.ts    # PWA service worker
  └── app.css              # Global styles

static/                    # Static assets
  ├── manifest.json
  └── icons/
```

### State Management

- **Svelte Stores:** Global state (auth, user preferences)
- **Component State:** Local UI state with `$state` runes
- **API Cache:** Service worker for offline support

---

## Database Schema

### Entity Relationship Diagram

```
┌──────────────┐
│    users     │
├──────────────┤
│ id (PK)      │
│ callsign     │◄───┐
│ email        │    │
│ tier         │    │
└──────────────┘    │
                    │
                    │
┌──────────────┐    │
│user_settings │    │
├──────────────┤    │
│ user_id (FK) │────┘
│ lotw_enabled │
│ grid_square  │
└──────────────┘


┌──────────────┐
│     qsos     │  (Hypertable)
├──────────────┤
│ id (PK)      │
│ user_id (FK) │────┐
│ callsign     │    │
│ frequency    │    │
│ band         │    │
│ mode         │    │
│ time_on      │◄───┼─── Partitioned by time
│ lotw_sent    │    │
└──────────────┘    │
                    │
┌──────────────────┐│
│propagation_data  ││  (Hypertable)
├──────────────────┤│
│ id (PK)          ││
│ timestamp        │◄─── Partitioned by time
│ solar_flux       │
│ k_index          │
└──────────────────┘
```

### Time-Series Optimization

- **QSOs Table:** Partitioned by `time_on` for efficient time-range queries
- **Propagation Data:** Partitioned by `timestamp` for historical analysis
- **Indexes:** Optimized for common queries (callsign, band, mode, date)

---

## Security Architecture

### Authentication Flow

```
User → Frontend → Backend → QRZ OAuth
                      ↓
                  Generate JWT
                      ↓
              Store in localStorage
                      ↓
        Include in all API requests
        (Authorization: Bearer {token})
```

### Security Measures

1. **Authentication:**
   - OAuth 2.0 via QRZ.com
   - JWT tokens with expiration
   - HTTP-only cookies (future enhancement)

2. **Authorization:**
   - Role-based access control (RBAC)
   - Tier-based feature gating

3. **Data Protection:**
   - HTTPS only (enforced in production)
   - CORS configuration
   - SQL injection prevention (parameterized queries)
   - XSS protection (CSP headers)

4. **Rate Limiting:**
   - Per-tier limits
   - API key throttling
   - DDoS protection (future: Cloudflare)

---

## API Design

### RESTful Principles

- **Versioned:** `/api/v1/...`
- **Resource-based:** `/api/v1/qso`, `/api/v1/user`
- **HTTP Methods:** GET, POST, PUT, DELETE
- **Status Codes:** Proper use of 2xx, 4xx, 5xx
- **Pagination:** Offset/limit for large datasets

### Future: GraphQL

GraphQL endpoint planned for complex queries:
```
POST /api/graphql
```

---

## Deployment Architecture

### Development

```
┌─────────────────────────────────────┐
│      Docker Compose (Local)         │
│  ┌─────────┐  ┌─────────┐          │
│  │Postgres │  │ Backend │           │
│  └─────────┘  └─────────┘           │
│       Frontend (Vite dev server)    │
└─────────────────────────────────────┘
```

### Production (Future)

```
┌─────────────────────────────────────────────┐
│                 Fly.io / PaaS               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │Frontend  │  │ Backend  │  │Postgres  │ │
│  │(Vercel)  │  │(Fly.io)  │  │(Supabase)│ │
│  └──────────┘  └──────────┘  └──────────┘ │
│                                             │
│  + Cloudflare CDN                          │
│  + Redis Cache (future)                    │
└─────────────────────────────────────────────┘
```

---

## Performance Considerations

### Backend
- Connection pooling for database
- Caching layer (Redis planned)
- Efficient SQL queries with indexes
- Goroutines for concurrent operations

### Frontend
- Code splitting (SvelteKit automatic)
- Lazy loading for routes
- Image optimization
- Service worker caching

### Database
- TimescaleDB compression for old data
- Partitioning by time
- Materialized views for analytics (future)

---

## Scalability

### Horizontal Scaling
- Stateless backend (JWT auth)
- Load balancer ready
- Database connection pooling

### Vertical Scaling
- Optimized queries
- Efficient data structures
- Minimal dependencies

---

## Monitoring & Observability (Future)

- **Logging:** Structured JSON logs
- **Metrics:** Prometheus + Grafana
- **Tracing:** OpenTelemetry
- **Error Tracking:** Sentry
- **Uptime Monitoring:** Better Uptime

---

## Technology Choices & Rationale

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Backend Language | Go | Fast, concurrent, simple deployment |
| Backend Framework | Fiber | Low latency, Express-like API |
| Database | PostgreSQL | Reliability, ACID, ecosystem |
| Time-Series | TimescaleDB | Optimized for QSO logs, solar data |
| Frontend | SvelteKit | Fast, modern, great DX, PWA support |
| Styling | Tailwind | Rapid development, consistent design |
| PWA | @vite-pwa | Offline support, mobile install |
| Container | Docker | Consistent environments, easy deployment |

---

*Last Updated: 2025-11-23*
