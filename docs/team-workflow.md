# Team Collaboration & Git Workflow

## 1. Branch Ownership by Team Member

After the Phase 0 foundation is merged to `main`, each member works in their designated feature branch:

| Member | Branch | Domain Focus | Initial GitHub Issue / Assignment |
| :--- | :--- | :--- | :--- |
| **Member 1 (Team Lead)** | `feat/subscription-lifecycle` | Subscriptions & Audit | Implement subscription CRUD, state machine transitions, renewal policies, and immutable `subscription_events` logging. |
| **Member 2** | `feat/catalog-plans` | Products & Pricing | Implement products, variant adjustments (e.g. Odoo connector +₹560), frequency plans, and administrative discount rules. |
| **Member 3** | `feat/billing-engine` | Recurring Billing & Cash | Implement recurring billing batch runner, idempotent invoice generation, explainable pricing preview, and payment recording. |
| **Member 4** | `feat/quotations-portal` | Quotations & Customer Portal | Implement quotation templates, proposal-to-subscription conversion, customer self-service portal, and executive revenue dashboard. |

---

## 2. Git Workflow Standards

1. **Keep `main` Green**: Direct commits to `main` are restricted. All changes enter via Pull Requests passing CI.
2. **Branch Naming**:
   - `feat/<feature-name>`
   - `fix/<bug-name>`
   - `docs/<topic>`
3. **Commit Messages**: Follow Conventional Commits:
   - `feat: add subscription pause policy check`
   - `fix: prevent duplicate billing on leap year end of month`
   - `test: add unit tests for discount rule limits`
   - `docs: update API contract for invoice PDF generation`
4. **Pull Request Protocol**:
   - Use the provided PR template (`.github/PULL_REQUEST_TEMPLATE.md`).
   - Every PR must verify that:
     1. Go tests pass (`go test -v ./...`).
     2. Code is formatted (`gofmt`).
     3. Frontend builds without TypeScript errors (`npm run build`).
     4. No SQL placed in HTTP handlers.
     5. Financial values use integer minor units.
   - Require at least 1 peer review approval before merging.
