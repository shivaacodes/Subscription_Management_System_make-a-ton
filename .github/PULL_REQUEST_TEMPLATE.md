## Description
<!-- Provide a clear summary of the changes and the module/issue addressed. -->

## Member & Feature Branch
- **Author:** Member [1 / 2 / 3 / 4]
- **Branch:** `feat/[subscription-lifecycle | catalog-plans | billing-engine | quotations-portal]`
- **Related Issue:** #<!-- issue number -->

## Type of Change
- [ ] `feat`: New feature or business capability
- [ ] `fix`: Bug fix
- [ ] `docs`: Documentation updates
- [ ] `test`: New or updated tests
- [ ] `refactor`: Code refactoring without behavior change
- [ ] `chore`: Build or dependency updates

## Architecture & Module Boundary Verification
- [ ] No SQL queries placed directly inside HTTP handlers.
- [ ] No circular module dependencies introduced.
- [ ] Constructor-based dependency injection used.
- [ ] Exact monetary values represented in integer minor units (paise/cents).
- [ ] Financial operations wrapped in database transactions where appropriate.
- [ ] Changes do not silently mutate historical invoice price snapshots.

## Testing & Verification
<!-- Describe manual and automated tests executed. -->
- [ ] Unit tests added / updated and passing (`go test -v ./...`)
- [ ] Frontend builds cleanly (`npm run build`)
- [ ] Database migrations tested up and down cleanly

## Checklist
- [ ] My code adheres to the project's Go and React formatting and linting rules.
- [ ] I have updated relevant documentation in `/docs`.
- [ ] Peer review obtained from at least one team member.
