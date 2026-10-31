# Veyra (Subflow) Architecture Documentation

## 1. System Overview

**Veyra** is a unified subscription management and recurring billing platform tailored for the business model of Brightloop Software (selling core subscription software, variants, and optional add-ons).

The platform replaces manual spreadsheets, uncoordinated renewal tracking, disconnected invoice generation, and fragmented reporting.

## 2. Core Business Pipeline

```
Product -> Variant -> Plan -> Quotation -> Subscription -> Recurring Billing -> Invoice -> Payment -> Reporting
```

Every module connects logically to this pipeline.

## 3. High-Level Architecture (Modular Monolith)

The system is architected as a modular monolith:

```
[ React 18 + Vite Frontend (SPA) ]
             |
   REST API / JSON over HTTP
   (HttpOnly Cookie Sessions)
             |
             v
 [ Go HTTP Server (Chi Router) ]
             |
      [ HTTP Handlers ]
             |
    [ Domain Services ]
             |
  [ SQL-First Repositories ]
             |
             v
      [ PostgreSQL 17 ]
```

### 3.1 Architectural Principles
1. **No Business Logic in HTTP Handlers**: Handlers are solely responsible for decoding HTTP requests, validating request envelopes, invoking domain services, and mapping output to JSON envelopes.
2. **No SQL Inside Handlers**: All database operations are sequestered within repository implementations.
3. **No Circular Dependencies**: Clear unidirectional dependency graphs: `cmd/api -> handlers -> services -> repositories -> models`.
4. **Constructor-Based Dependency Injection**: Explicit struct initialization with dependencies passed via constructors (e.g. `auth.NewService(userRepo, sessionRepo)`).
5. **No Generic Repository Framework or Heavy ORM**: Using SQL-first queries with `pgx/v5` and `pgxpool`.
6. **UTC Timestamps**: All temporal fields stored and calculated in UTC (`TIMESTAMPTZ`).
7. **Exact Monetary Representation**: Currency stored in integer minor units (paise/cents as `BIGINT`) to eliminate floating-point arithmetic errors.
8. **Financial Transactions**: All financial state mutations (invoices, payments, subscriptions) execute within explicit database transactions (`tx.Begin` -> `tx.Commit`).
9. **Immutable Invoice Pricing Snapshots**: Confirmed invoices store snapshotted unit prices, discounts, and tax figures to preserve historical accounting records.
10. **Explainable Billing**: One authoritative pricing service powers both preview and generation: `Base Price -> Variant Adjustment -> Quantity -> Discount -> Tax -> Final Total`.

## 4. Domain Module Boundaries

- `auth`: User registration, session creation, password validation, role checks (`ADMIN`, `STAFF`, `CUSTOMER`).
- `users`: User identity and customer profiles.
- `catalog`: Products, types, variants (e.g., Odoo Brand connector +₹560), and pricing adjustments.
- `plans`: Frequency configurations (DAILY, WEEKLY, MONTHLY, YEARLY), availability dates, and policies (auto-close, allow pause, allow renewal).
- `quotations`: Quotation templates, proposals, expiration dates, line items, and subscription conversion.
- `subscriptions`: Active subscriptions, lines, status state machine (`DRAFT -> QUOTATION -> CONFIRMED -> ACTIVE -> CLOSED`), add-on toggling.
- `billing`: Recurring billing runs, idempotent invoice generation, Billing Reliability Center, line snapshots.
- `payments`: Payment records, multi-part invoice allocations, derived invoice payment states.
- `reporting`: MRR, invoiced revenue, cash collected, accounts receivable, and aging buckets.
- `audit`: Immutable event sourcing and history timeline for subscriptions (`subscription_events`).
