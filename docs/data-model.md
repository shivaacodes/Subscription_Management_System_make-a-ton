# Veyra Data Model Specification

## 1. Relational Entities & Ownership Boundaries

### Core Identity & Access Control (Foundation)
- `users`: Core identity table storing normalized email (`LOWER(email)` unique index), bcrypt password hash, and assigned role (`ADMIN`, `STAFF`, `CUSTOMER`).
- `customer_profiles`: Detailed company, telephone, and billing address profile linked 1:1 with users having the `CUSTOMER` role.
- `sessions`: Ephemeral token hashes for server-side HttpOnly session management, tracking user agent, IP address, and expiry.
- `password_reset_tokens`: Single-use expiring tokens for secure password recovery.

### Catalog & Pricing (Member 2 Ownership: `feat/catalog-plans`)
- `products`: Base sellable goods or software services. Stored with `sales_price_cents` and `cost_price_cents` as `BIGINT` minor units.
- `product_variants`: Configurable add-on pricing options (e.g. Odoo Brand connector adds ₹560 adjustment).
- `plans`: Frequency-based subscription definitions (`DAILY`, `WEEKLY`, `MONTHLY`, `YEARLY`) with min quantity, validity dates, and lifecycle policy flags (`auto_close`, `allow_early_close`, `allow_pause`, `allow_renewal`).
- `plan_products`: Join table declaring standard and optional add-on products for a plan.
- `tax_rules`: Tax jurisdictions and rates in basis points (`rate_bps`, e.g. 1800 for 18.00% GST).
- `discount_rules`: Fixed amount or percentage discount rules with minimum spend, max usage counters, and product/subscription targets.

### Quotations & Experience (Member 4 Ownership: `feat/quotations-portal`)
- `quotation_templates`: Reusable quotation templates with preconfigured lines.
- `quotation_template_lines`: Line items associated with templates.
- `quotations`: Formal client quotations with validity periods, subtotal, tax, and discount snapshots.
- `quotation_lines`: Itemized product/variant lines with applied tax rates.

### Subscriptions & Timeline (Member 1 / Team Lead Ownership: `feat/subscription-lifecycle`)
- `subscriptions`: Primary contract record between customer and plan. Unique human-readable `subscription_number`.
  - **State Machine**: `DRAFT -> QUOTATION -> CONFIRMED -> ACTIVE -> CLOSED` (with `PAUSED` and `CANCELLED` states).
  - Next billing date and payment terms (`NET_30`, etc.).
- `subscription_lines`: Contract line items capturing quantity, unit price, tax rate, and totals.
- `subscription_events`: Immutable audit event stream capturing all state transitions, pricing changes, and administrative actions.

### Recurring Billing & Payments (Member 3 Ownership: `feat/billing-engine`)
- `billing_runs`: Audit records of automated batch recurring billing runs.
- `billing_attempts`: Log of individual subscription billing attempts with retry status.
- `invoices`: Formal billing statement with unique `invoice_number` and composite uniqueness:
  `UNIQUE (subscription_id, billing_period_start, billing_period_end)` guarantees database-level idempotency against duplicate billing.
- `invoice_lines`: Frozen pricing snapshot of items billed in this period.
- `discount_applications`: Record of discount rule applied to an invoice.
- `payments`: Multi-part cash/bank transfers applied to invoices with unique `payment_number`.

---

## 2. Monetary & Numeric Conventions

1. **Integer Minor Units**: All monetary values are represented as `BIGINT` representing cents/paise (e.g., ₹4,999.00 = `499900`).
2. **Basis Points for Percentages**: Tax and percentage discount values are stored in basis points (1 bp = 0.01%, 1800 bps = 18.00%).
3. **No Floating Point Operations**: Floating-point numbers (`float64`, `REAL`) are forbidden for monetary arithmetic in Go and PostgreSQL.

---

## 3. Idempotency & Concurrency Model

- **Duplicate Billing Prevention**: The database constraint `uq_sub_billing_period` on `(subscription_id, billing_period_start, billing_period_end)` ensures that no two invoices can be created for the same billing window, even under concurrent execution.
- **Transactional Consistency**: Invoice creation, line item insertion, and subscription `next_billing_date` increments occur in a single database transaction.
