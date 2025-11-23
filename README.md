## Ham‑Radio Cloud – Product Blueprint (v1.1 - Integrated Suggestions)

### 1 Executive Snapshot

A subscription SaaS + API platform that unifies **station logging, real‑time propagation insights, SDR‑stream hosting, and advanced DSP tooling** for the global amateur‑radio community. The service off‑loads database headaches for operators, provides one‑click logbook‑of‑the‑world (LoTW) sync, and exposes REST/WebSocket endpoints so hardware makers can embed cloud analytics.

### 2 Problem & Opportunity

- Hams juggle aging Windows logbook apps, separate solar‑data sites, DIY SDR servers, and CSV exports to ARRL/LoTW.
- Mobile & remote‑shack operation is growing (field days, POTA/SOTA). They crave **cloud sync + phone/tablet dashboards**.
- The niche is passionate, spends on gear, and has light regulatory friction (FCC Part 97 allows remote control with safeguards).

### 3 Core Value Proposition

| Stakeholder          | Pain                                        | Our Fix                                                        |
| -------------------- | ------------------------------------------- | -------------------------------------------------------------- |
| **Casual Operator**  | Manual log entries, lost QSOs, LoTW hassle  | Auto‑log via CAT/ADIF (WebUSB/rigctld), cloud backup, 1-click LoTW |
| **Contester / DXer** | Need propagation intel & cluster alerts     | Integrated solar/Kp data + VOACAP baseline + live DX cluster push |
| **SDR Enthusiast**   | Remote access complexity, sharing difficulty| Simplified public stream integration; private hosting post-MVP |
| **Hardware Makers**  | No cloud layer                              | White‑label API to display stats/log QSOs in their apps         |

### 4 Product Modules

1.  **Cloud Logbook (Core)**
    *   ADIF import/export, LoTW/eQSL auto‑sync, award tracker.
    *   Manual QSO entry, basic filtering/reporting.
    *   Support for WebUSB & local `rigctld` integration via helper app/agent.
2.  **Propagation Engine**
    *   Ingest solar‑flux (SIDC), Kp; integrate baseline VOACAP predictions.
    *   Rule-based band opening alerts (MVP).
    *   *Future:* LSTM model predicts band openings per grid square (post-MVP, experimental).
3.  **Web‑SDR Integration**
    *   *MVP:* Integration with existing *public* WebSDR directories (e.g., KiwiSDR network).
    *   *Future:* Docker‑based user-hosted SDR receiver pod (read-only first).
    *   *Future:* Optional transmit capability (requires strict controls, user callsign whitelist, PTT watchdog, enhanced auth).
4.  **DSP Toolbox (Future)**
    *   Browser‑side WASM FIR filters, noise reduction, CW/PSK decoder (Post-MVP).
5.  **REST / GraphQL API & Webhooks**
    *   Spotting cluster feed, QSO write, propagation query, station tele‑metrics.
6.  **Admin & Billing**
    *   OAuth (QRZ.com / Google), Stripe metered seats/subscriptions, club billing groups.

### 5 Tech Stack

| Layer          | Choice                                         | Reason                                                |
| -------------- | ---------------------------------------------- | ----------------------------------------------------- |
| Backend        | Go + Fiber (low‑latency)                       | Concurrency for WebSockets, API efficiency            |
| DB             | Postgres (Timescale ext)                       | Time‑series for QSO data, solar stats, metrics        |
| AI (Future)    | Python micro‑service (FastAPI, PyTorch)        | Train & serve advanced propagation model (post-MVP)   |
| SDR Containers (Future) | Docker + Rust wrappers                 | Secure capability drop for user-hosted SDRs          |
| Front‑End      | SvelteKit + Tailwind                           | Lightweight, good realtime updates, PWA capable       |
| Hosting        | Fly.io or similar PaaS; Hetzner for SDR nodes  | Edge locations, manageable infra, cost efficiency     |

### 6 Compliance / Licensing

- **FCC Part 97 / CEPT etc.**: Remote TX requires authenticated control operator (OAuth + callsign + MFA suggested). Strict 15‑minute PTT watchdog, per-band power limits, geofencing per license/band plan. Clear user opt-in & warnings.
- Privacy: Minimal PII (callsigns mostly public). GDPR-ready consent banner, clear data policy.

### 7 MVP (Focus: ~60-90 days)

1.  **Rock-Solid Cloud Logbook:**
    *   Manual QSO Entry & ADIF Import/Export.
    *   Reliable LoTW Sync (one-way push initially acceptable).
    *   Basic QSO filtering/viewing.
2.  **Basic Propagation Display:**
    *   Ingest Solar/Kp data.
    *   Display current conditions & basic rule-based "Good/Fair/Poor" band indicators (VOACAP integration if feasible).
3.  **Public Web-SDR Directory Integration:**
    *   List/link to existing public KiwiSDRs (read-only). No user hosting yet.
4.  **Core Platform:**
    *   User Auth (OAuth - QRZ.com essential).
    *   Mobile-first PWA Dashboard (SvelteKit).
    *   Stripe Subscription Setup (Free & Operator Tiers initially).

### 8 Pricing Model (Simplified Initial Tiers)

| Tier        | Features                                                   | Monthly |
| ----------- | ---------------------------------------------------------- | ------- |
| Free        | 500 QSOs limit, basic propagation display, public SDR links| \$0     |
| Operator    | 20k QSOs, LoTW sync, rule-based alerts, club groups        | \$7     |
| *Future:* Contester | Unlimited QSOs, priority features, private SDR hosting, API access | \$19+   |
| *Future:* Partner API | Usage-based API access, OEM license                | \$99+   |

*Note: Remote TX and advanced SDR/API features will be add-ons or part of higher tiers introduced post-MVP.*

### 9 Go‑to‑Market

1.  **Seed with Influencers**: Sponsor YouTube ham channels (Dave Casler, HRC, etc.). Provide early access.
2.  Launch on **r/amateurradio**, QRZ.com forums; offer beta invites for the "Operator" tier.
3.  Club outreach: ARRL sections, contest clubs – focus on the "Club Dashboard" potential and offer group trials/discounts.
4.  Booth at **Dayton Hamvention** – demo the polished logging & LoTW sync, preview roadmap.

### 10 12‑Month Roadmap (Post-MVP)

| Q  | Theme                                       | Key Features                                                      |
| -- | ------------------------------------------- | ----------------------------------------------------------------- |
| Q1 | Stabilize MVP, User Feedback, Scale Basics | Performance tuning, UI refinements, expand LoTW sync capabilities |
| Q2 | Enhance Propagation & Introduce Private SDR | VOACAP integration deep-dive, Beta AI model, Private SDR (RX only)|
| Q3 | Monetization Expansion & Remote TX Beta   | Contester Tier launch, API Beta, Remote TX (controlled beta)       |
| Q4 | Hardware Integration & Community Features   | rigctld agent polish, Partner SDK v1, Advanced award tracking     |

### 11 Risk & Mitigation

| Risk                      | Impact | Mitigation                                                               |
| ------------------------- | ------ | ------------------------------------------------------------------------ |
| **MVP Scope Creep**       | High   | Ruthless focus on core logging/LoTW sync reliability for initial launch.   |
| **Remote TX Misuse**      | High   | Strict auth (MFA), per-band power caps, callsign logging, PTT watchdog, phased rollout. |
| **SDR Hosting Costs/Perf**| Med    | Delay user-hosting; start with public links; monitor infra costs closely. |
| **AI Forecast Accuracy**  | Med    | Start with baseline (VOACAP/rules); gradual AI rollout; manage user expectations. |
| Niche TAM Ceiling         | Med    | Future: Expand to marine HF, SWL markets, adjacent radio hobbies.      |

### 12 North‑Star Metrics

- Weekly Active Loggers (WAL) - *Focus on core logging activity first.*
- LoTW Sync Success Rate (>99%)
- Mean QSO logging latency (<1.5s via UI/API)
- *Future:* Forecast hit‑rate (>70% validated band‑open precision)
- Churn < 3% / mo after 90 days (post-paid launch)

### 13 Phased Roadmap to MVP

To ensure a focused, iterative path to MVP, development will proceed in the following phases. Each phase builds on the previous, enabling early feedback and rapid delivery of core value.

### Phase 1: Project Foundation & Core Infrastructure
- Initialize monorepo (backend, frontend, shared).
- Set up CI/CD, linting, Prettier, and basic test harness.
- Provision Postgres (with Timescale extension) and basic schema for users, QSOs, and logs.
- Deploy basic backend skeleton (Go + Fiber) with health checks.
- Set up SvelteKit frontend with Tailwind, PWA config, and initial routing.
- Integrate OAuth (QRZ.com) for user authentication (backend + frontend).
- **Acceptance:** User can sign up/log in, see a dashboard placeholder, and backend is deployed with DB connectivity.

### Phase 2: Cloud Logbook Core
- Manual QSO entry UI and backend API.
- ADIF import/export (UI + backend).
- QSO list view with basic filtering (date, callsign, band).
- LoTW sync (one-way push) integration (backend job + UI status).
- Basic error handling and user feedback.
- **Acceptance:** User can log QSOs, import/export ADIF, and push logs to LoTW.

### Phase 3: Propagation Engine (Basic)
- Ingest solar/Kp data (scheduled backend job).
- Store and expose propagation data via API.
- Display current solar/Kp conditions in dashboard.
- Rule-based "Good/Fair/Poor" band indicators (VOACAP baseline if feasible).
- **Acceptance:** User sees real-time propagation data and band conditions on dashboard.

### Phase 4: Public Web-SDR Directory Integration
- Integrate and display public KiwiSDR directory (read-only).
- UI component for browsing/listing SDR streams.
- Link out to public SDRs (no user hosting yet).
- **Acceptance:** User can browse and access public SDR streams from the dashboard.

### Phase 5: Core Platform & Billing
- Stripe integration for subscriptions (Free & Operator tiers).
- Enforce QSO limits and feature gating by tier.
- Admin dashboard for user management and basic metrics.
- Mobile-first PWA polish (install prompt, offline fallback).
- **Acceptance:** Users can upgrade/downgrade, limits enforced, and app is PWA-ready.

### Phase 6: Testing, Hardening, and Launch Prep
- End-to-end and integration tests for critical flows.
- Security review (OAuth, data privacy, rate limiting).
- Performance tuning (API latency, DB queries).
- Documentation (README, API docs, onboarding guide).
- Beta invite system and feedback loop.
- **Acceptance:** All tests pass, security checks complete, and ready for controlled beta.

---

*Each phase should be tracked in `tasks/tasks.md` and progress/status in `docs/status.md`. Compliance and security are considered from the start, especially for user data and LoTW integration. After MVP, the roadmap can expand to private SDR hosting, advanced DSP, and API/partner features as outlined above.*

---

*v1.1 - Ready for initial project setup.* 
