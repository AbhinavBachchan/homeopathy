# Homeopathy Platform — Go + Angular Starter

This is a Phase 1 (P0 launch scope) starter scaffold adapted from the developer
brief, swapping Next.js/Node for **Go (Gin + GORM)** on the backend and
**Angular 18 (standalone components)** on the frontend.

## What's included

**Backend (`/backend`)**
- Gin REST API with JWT auth (register/login), role-based middleware
  (patient/doctor/admin/corporate_hr)
- GORM models for `users`, `products`, `orders` mapped from the brief's schema
  (section 9.1), including homeopathy-specific product fields (potency, form,
  schedule, indications, contraindications)
- Product listing with filters (potency/brand/category/search) and admin CRUD
- Order creation with GST calculation and Schedule H prescription enforcement
- Config loader pre-wired for every third-party key in the brief's env var
  list (Razorpay, Stripe, Claude, Interakt, Brevo, MSG91, Shiprocket, Algolia)
  — the fields exist, the actual API calls are marked as `TODO`

**Frontend (`/frontend`)**
- Angular 18 standalone app (no NgModules) with routing, HTTP client + JWT
  interceptor
- Product list + product detail pages wired to the Go API
- Login page with reactive forms
- Signal-based `AuthService` for current-user state

**Infra**
- `docker-compose.yml` running Postgres, Redis, backend, and frontend together
- Dockerfiles for both services

## What's deliberately stubbed (next to build)

These are marked with `TODO` comments in the code, matching the brief's P0/P1
priority order:
- Mobile OTP login (MSG91) and Google OAuth
- Razorpay/Stripe payment intent creation + webhook handlers
- Shiprocket shipping rate calculation and label generation
- WhatsApp order notifications (Interakt)
- Bulk CSV product import
- Cart component, checkout flow, address management on the frontend
- Admin dashboard UI

Phase 2 (AI consultation) and Phase 3 (subscriptions/corporate) tables and
routes aren't scaffolded yet — add them once Phase 1 is live, per the brief's
"do not delay Phase 1" guidance.

## Running locally

**Prereqs:** Go 1.22+, Node 20+, Docker (optional but easiest)

### Option A — Docker Compose (recommended for a clean start)
```bash
docker compose up --build
```
- Backend: http://localhost:8080
- Frontend: http://localhost:4200
- Postgres: localhost:5432

### Option B — Run natively
```bash
# Backend
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/api

# Frontend (separate terminal)
cd frontend
npm install
npm start
```

The backend auto-migrates the database on startup (`db.AutoMigrate`), so once
Postgres is up, tables are created automatically — no separate migration step
needed yet. Swap in `golang-migrate` for versioned migrations once the schema
stabilizes.

## API quick reference

| Method | Endpoint | Auth |
|---|---|---|
| POST | `/api/auth/register` | Public |
| POST | `/api/auth/login` | Public |
| GET | `/api/products?potency=&brand=&category=&q=` | Public |
| GET | `/api/products/:slug` | Public |
| POST | `/api/orders` | Patient |
| GET | `/api/orders` | Patient (own orders) |
| GET | `/api/orders/:id` | Patient |
| POST | `/api/admin/products` | Admin |
| PUT | `/api/admin/products/:id` | Admin |
| DELETE | `/api/admin/products/:id` | Admin |

## Suggested next steps (matches brief Sprint 1–4)

1. Wire Razorpay checkout (`/api/webhooks/razorpay` + client-side Razorpay
   Checkout.js in Angular)
2. Build the cart feature (Angular service holding cart state, `/orders`
   integration)
3. Add Shiprocket integration for rates + label generation
4. Add admin product CRUD screens in Angular (table + form)
5. Bulk CSV import endpoint (`POST /api/admin/products/import`)
6. Interakt WhatsApp order-status webhook triggers

This gets you to "B2C Store fully live" (brief's Sprint 4 milestone) without
Phase 2/3 features blocking the launch.
