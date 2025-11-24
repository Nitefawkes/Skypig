# Claude Development Guide - Ham-Radio Cloud (Skypig)

## Project Overview

Ham-Radio Cloud is a subscription SaaS + API platform that unifies **station logging, real-time propagation insights, SDR-stream hosting, and advanced DSP tooling** for the global amateur-radio community. The platform off-loads database management for operators, provides one-click Logbook-of-the-World (LoTW) sync, and exposes REST/WebSocket endpoints for hardware integration.

**Current Status:** v1.1 - Project initialization and foundation phase
**Target MVP:** 60-90 days

## Core Value Propositions

- **Cloud Logbook:** Auto-logging via CAT/ADIF, cloud backup, 1-click LoTW sync
- **Propagation Engine:** Real-time solar/Kp data with band opening alerts
- **Web-SDR Integration:** Simplified access to public SDR streams (future: private hosting)
- **API Platform:** White-label API for hardware makers and third-party integrations

## Tech Stack

### Backend
- **Language:** Go
- **Framework:** Fiber (low-latency, optimized for WebSockets)
- **Database:** PostgreSQL with Timescale extension
  - Time-series data for QSO logs, solar stats, metrics
- **Future:** Python microservice (FastAPI, PyTorch) for AI propagation models

### Frontend
- **Framework:** SvelteKit
- **Styling:** Tailwind CSS
- **PWA:** Mobile-first, progressive web app capable
- **Real-time:** WebSocket support for live updates

### Infrastructure
- **Hosting:** Fly.io or similar PaaS
- **SDR Nodes (Future):** Hetzner for dedicated SDR infrastructure
- **CI/CD:** GitHub Actions (to be configured)

### Future Modules
- **SDR Containers:** Docker + Rust wrappers with security capability drops
- **DSP Toolbox:** Browser-side WASM for FIR filters, noise reduction, CW/PSK decoder

## Development Phases

We're following a phased approach to MVP:

### Phase 1: Project Foundation & Core Infrastructure (CURRENT)
- Initialize monorepo (backend, frontend, shared)
- Set up CI/CD, linting, Prettier, and test harness
- Provision Postgres with Timescale extension
- Deploy basic backend skeleton (Go + Fiber) with health checks
- Set up SvelteKit frontend with Tailwind and PWA config
- Integrate OAuth (QRZ.com authentication)

### Phase 2: Cloud Logbook Core
- Manual QSO entry UI and API
- ADIF import/export functionality
- QSO list view with filtering
- LoTW sync (one-way push)

### Phase 3: Propagation Engine (Basic)
- Ingest solar/Kp data via scheduled jobs
- Display current conditions
- Rule-based band indicators (Good/Fair/Poor)
- VOACAP baseline integration

### Phase 4: Public Web-SDR Directory Integration
- Integrate KiwiSDR directory
- UI for browsing public SDR streams
- Read-only access (no user hosting yet)

### Phase 5: Core Platform & Billing
- Stripe subscription integration
- Feature gating by tier (Free/Operator)
- Admin dashboard
- PWA polish and offline support

### Phase 6: Testing, Hardening, and Launch Prep
- E2E and integration tests
- Security review
- Performance tuning
- Beta invite system

## Code Organization

```
/
├── backend/          # Go + Fiber API server
├── frontend/         # SvelteKit PWA
├── shared/           # Shared types, utilities
├── tasks/            # Task tracking (tasks.md)
├── docs/             # Documentation (status.md, etc.)
├── README.md         # Product blueprint
└── claude.md         # This file
```

## Development Guidelines

### General Principles

1. **MVP Focus:** Ruthlessly prioritize core features. Avoid scope creep.
2. **Security First:** Always consider FCC Part 97 compliance, authentication, and data privacy.
3. **Mobile-First:** Design and test for mobile/tablet usage patterns.
4. **Performance:** Target <1.5s QSO logging latency via UI/API.
5. **Reliability:** Aim for >99% LoTW sync success rate.

### Coding Standards

#### Backend (Go)
- Follow standard Go conventions and idioms
- Use meaningful variable and function names
- Leverage Go's concurrency primitives for WebSocket handling
- Implement proper error handling and logging
- Write table-driven tests for business logic
- Use Fiber's middleware for auth, CORS, rate limiting

#### Frontend (SvelteKit)
- Component-based architecture
- Use Tailwind utility classes, avoid custom CSS where possible
- Implement reactive stores for state management
- Progressive enhancement: ensure core features work without JS
- Optimize bundle size and lazy-load routes
- Use TypeScript for type safety

#### Database
- Use migrations for schema changes (track in version control)
- Leverage Timescale for time-series queries
- Index appropriately for common query patterns (callsign, date ranges)
- Normalize user data but denormalize for read-heavy QSO queries where beneficial

### Security Considerations

- **Authentication:** OAuth via QRZ.com (essential for ham radio community)
- **Authorization:** Implement role-based access control (RBAC)
- **Data Privacy:** Minimal PII collection, GDPR-ready consent
- **API Security:** Rate limiting, API key management for partners
- **Future Remote TX:** MFA, callsign verification, PTT watchdog, geofencing, power limits

### Testing Strategy

- **Unit Tests:** Go backend business logic, critical frontend utilities
- **Integration Tests:** API endpoints, database operations
- **E2E Tests:** Critical user flows (login, QSO entry, LoTW sync)
- **Manual Testing:** Mobile responsiveness, PWA install, cross-browser

### API Design

- **REST:** Standard CRUD operations for QSOs, users, settings
- **WebSocket:** Live updates for DX cluster, propagation alerts
- **GraphQL (Future):** Complex queries for reporting and analytics
- **Webhooks:** Event notifications for partner integrations

### Compliance & Regulatory

- **FCC Part 97:** Remote control operator authentication, PTT watchdog (15min), power limits, band plans
- **CEPT/International:** Geofencing per license/band plan
- **Data Retention:** Clear policies for QSO data, callsign information
- **User Consent:** Explicit opt-in for remote TX features

## Key Integrations

### Current/MVP
- **QRZ.com OAuth:** User authentication
- **LoTW (ARRL):** One-way push for QSO log submission
- **SIDC Solar Data:** Real-time solar flux and Kp index
- **KiwiSDR Directory:** Public Web-SDR access

### Future
- **VOACAP:** HF propagation prediction
- **rigctld:** Local rig control daemon integration
- **eQSL:** Alternative QSL card service
- **Contest APIs:** Integration with major contest platforms

## Pricing Tiers

| Tier | Features | Monthly |
|------|----------|---------|
| Free | 500 QSOs, basic propagation, public SDR links | $0 |
| Operator | 20k QSOs, LoTW sync, alerts, club groups | $7 |
| Contester (Future) | Unlimited QSOs, private SDR, API access | $19+ |
| Partner API (Future) | Usage-based API, OEM license | $99+ |

## Metrics to Track

- **Weekly Active Loggers (WAL)** - Primary engagement metric
- **LoTW Sync Success Rate** - Target >99%
- **Mean QSO Logging Latency** - Target <1.5s
- **Churn Rate** - Target <3% monthly after 90 days
- **Future: Forecast Hit Rate** - Target >70% for band opening predictions

## Common Tasks

### Starting Development
```bash
# Backend
cd backend
go mod download
go run main.go

# Frontend
cd frontend
npm install
npm run dev
```

### Database Setup
```bash
# Initialize Postgres with Timescale
docker-compose up -d postgres
# Run migrations
cd backend && go run migrations/migrate.go
```

### Running Tests
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm test
```

## Resources & References

- **FCC Part 97:** https://www.ecfr.gov/current/title-47/chapter-I/subchapter-D/part-97
- **LoTW Documentation:** https://lotw.arrl.org/lotw-help/
- **ADIF Specification:** https://adif.org/
- **KiwiSDR:** http://kiwisdr.com/
- **VOACAP:** https://www.voacap.com/

## Risk Mitigation

- **MVP Scope Creep:** Focus on core logging/LoTW sync reliability first
- **Remote TX Misuse:** Strict auth, logging, watchdog, phased rollout
- **SDR Hosting Costs:** Start with public links, delay user hosting
- **AI Accuracy:** Begin with VOACAP/rules, manage expectations

## Communication Style

When working with this codebase:
- Prioritize reliability over features
- Consider mobile/field usage scenarios
- Respect amateur radio community norms and terminology
- Think about club and group use cases
- Plan for international users (CEPT, licensing variations)

## Questions to Consider

Before implementing features:
- Does this serve the MVP goals?
- What's the mobile experience?
- How does this scale for clubs/groups?
- Are there regulatory implications?
- What's the cost/infrastructure impact?
- Does it integrate with existing ham radio workflows?

---

*Last Updated: 2025-11-24*
*Version: 1.0 - Initial comprehensive guide*
