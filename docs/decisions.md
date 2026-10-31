# Architecture Decision Records (ADRs)

## ADR-001: Modular Monolith vs Microservices
- **Status:** Accepted
- **Context:** Four-member team building an MVP for an Odoo Hackathon evaluation within tight deadlines.
- **Decision:** Build a single Go modular monolith with clear package boundaries (`Handler -> Service -> Repository`).
- **Consequences:** Avoids distributed tracing, network latency, multiple deployment pipelines, and serialization overhead while maintaining clean modular isolation.

## ADR-002: Go Standard Library & Chi Router over Frameworks
- **Status:** Accepted
- **Context:** Standardized Go HTTP stack with minimal dependencies.
- **Decision:** Use standard `net/http` server with `go-chi/chi/v5` router.
- **Consequences:** Idiomatic Go code, low memory footprint, zero framework lock-in, and full compatibility with standard Go middleware and testing tools.

## ADR-003: Raw SQL with pgx/v5 over GORM or Other ORMs
- **Status:** Accepted
- **Context:** Complex financial calculations, composite unique constraints, and strict transactional control.
- **Decision:** Use raw SQL queries with `jackc/pgx/v5` and connection pooling (`pgxpool`).
- **Consequences:** Clear understanding of generated queries, direct control over transactions, zero ORM caching bugs, and precise constraint mapping.

## ADR-004: Server-Side HttpOnly Session Cookies over JWT in localStorage
- **Status:** Accepted
- **Context:** Secure web session management.
- **Decision:** Generate 32-byte cryptographically secure random session tokens stored as SHA-256 hashes in PostgreSQL `sessions` table. Transmit via `HttpOnly`, `SameSite=Lax`, and `Secure` cookies.
- **Consequences:** Immune to client-side token theft via XSS; immediate server-side revocation on logout or account compromise; no stale JWT expiration dilemmas.

## ADR-005: Integer Minor Units for Monetary Values
- **Status:** Accepted
- **Context:** Invoicing, taxes, line totals, and payments require exact arithmetic.
- **Decision:** Represent currency as `BIGINT` in minor units (cents / paise).
- **Consequences:** Completely prevents floating-point rounding errors (e.g. `0.1 + 0.2 != 0.3`). Simple to render in frontend via `formatCurrency(cents / 100)`.

## ADR-006: Composite Unique Constraint for Idempotent Billing
- **Status:** Accepted
- **Context:** Recurring billing jobs run periodically and must never generate duplicate invoices.
- **Decision:** Define `UNIQUE (subscription_id, billing_period_start, billing_period_end)` on the `invoices` table.
- **Consequences:** Concurrency-safe duplicate prevention enforced directly by the database engine.
