# Veyra

> Unified Subscription Management and Recurring Billing Platform for Brightloop Software.  
> Built for the Odoo Hackathon | Phase 0 Technical Foundation.

---

## Architecture Summary

Veyra is structured as a high-reliability **modular monolith** with strict layer separation:

```
[ Frontend (React 18 + Vite + TS) ]  <--- HttpOnly Secure Cookies --->  [ Backend (Go 1.24+ + Chi) ]
                                                                                   |
                                                                        [ Repositories (pgx/v5) ]
                                                                                   |
                                                                        [ PostgreSQL 17 Engine ]
```

### Key Pillars
- **Strict Layering**: `Handler -> Service -> Repository`. No business logic or SQL inside HTTP handlers.
- **Explainable Billing**: Autoritative pricing service preview: `Base Price -> Variant Adjustment -> Quantity -> Discount -> Tax -> Final Total`.
- **Idempotent Recurring Billing**: Composite database constraint `UNIQUE(subscription_id, billing_period_start, billing_period_end)` prevents duplicate invoicing.
- **Accurate Financials**: Exact integer minor unit monetary modeling (`BIGINT` paise/cents). No floating-point inaccuracies.
- **Security**: HttpOnly, SameSite=Lax session cookies, bcrypt password hashing, and strict password policy (>8 chars, uppercase, lowercase, special character).

---

## Repository Structure

```text
.
├── frontend/                     # React 18 + TypeScript SPA
│   ├── src/
│   │   ├── app/
│   │   ├── components/ui/       # Restrained UI components (Button, Card, Input, Alert)
│   │   ├── components/layout/   # Public, Admin/Staff, and Customer Portal layouts
│   │   ├── features/            # Auth, Subscriptions, Catalog, Billing, Quotations
│   │   ├── lib/api/             # Native fetch API client with credentials support
│   │   ├── routes/              # Route definitions with role-based protection
│   │   └── types/               # TypeScript interfaces
│   ├── package.json
│   └── vite.config.ts
├── backend/                      # Go Modular Monolith API
│   ├── cmd/
│   │   ├── api/                 # API server entrypoint (port 8080)
│   │   └── seed/                # Standalone database seed utility
│   ├── internal/
│   │   ├── config/              # Environment config parsing
│   │   ├── platform/            # DB pool, migrator, json helpers, auth middleware
│   │   ├── auth/                # Auth service, session management, password policy
│   │   ├── users/               # User repository & customer profiles
│   │   ├── catalog/             # Products and variant contracts
│   │   ├── plans/               # Subscription plan contracts
│   │   ├── quotations/          # Quotation & template contracts
│   │   ├── subscriptions/       # Subscription lifecycle & audit timeline contracts
│   │   ├── billing/             # Billing run, invoice & calculation contracts
│   │   ├── payments/            # Payment recording contracts
│   │   ├── reporting/           # Revenue metrics contracts
│   │   └── audit/               # Timeline event models
│   ├── migrations/              # Versioned SQL migrations
│   └── go.mod
├── docs/                         # Technical documentation
│   ├── architecture.md
│   ├── data-model.md
│   ├── api-contracts.md
│   ├── decisions.md
│   ├── development.md
│   └── team-workflow.md
├── .github/
│   ├── workflows/ci.yml         # GitHub Actions pipeline
│   └── PULL_REQUEST_TEMPLATE.md
├── docker-compose.yml            # PostgreSQL 17 container
└── README.md
```

---

## Quickstart

### 1. Database Setup
Start PostgreSQL via Docker:
```bash
docker-compose up -d postgres
```

### 2. Configuration & Seeding
```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Run migrations and seed test credentials
cd backend
go run cmd/seed/main.go
```

**Pre-seeded Accounts**:
- **Admin**: `admin@brightloop.local` / `AdminPassword123!`
- **Staff**: `staff@brightloop.local` / `StaffPassword123!`
- **Customer**: `customer@brightloop.local` / `CustomerPassword123!`

### 3. Start Backend
```bash
cd backend
go run cmd/api/main.go
```
API running on `http://localhost:8080`.  
Health check: `curl http://localhost:8080/api/v1/health`.

### 4. Start Frontend
```bash
cd frontend
npm install
npm run dev
```
Open `http://localhost:5173`.

---

## Four-Member Feature Branch Tracks

| Member | Branch | Responsibilities |
| :--- | :--- | :--- |
| **Member 1 (Lead)** | `feat/subscription-lifecycle` | Subscriptions CRUD, lifecycle state transitions, renewals, add-on toggling, audit timeline |
| **Member 2** | `feat/catalog-plans` | Products, variants (+₹560 Odoo brand), plans, taxes, admin discount rules |
| **Member 3** | `feat/billing-engine` | Recurring billing runner, idempotent invoicing, Billing Reliability Center, payments |
| **Member 4** | `feat/quotations-portal` | Quotation templates, proposals to subscriptions, Customer portal, revenue reports |

See `docs/team-workflow.md` for guidelines.
