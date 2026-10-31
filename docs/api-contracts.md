# Veyra API Contracts

## 1. Base URL & Protocol
- **Base Path**: `/api/v1`
- **Session Transport**: HttpOnly Secure Cookies (`subflow_session`)
- **Headers**:
  - `Content-Type: application/json`
  - `X-Requested-With: XMLHttpRequest` (for mutating requests)

---

## 2. Standard Response Envelope

### Success Response
Direct JSON object or array:
```json
{
  "user": {
    "id": "1299763b-31a2-45ca-a47d-f1bce74290c0",
    "email": "user@example.com",
    "role": "CUSTOMER",
    "createdAt": "2026-10-31T12:55:18Z",
    "updatedAt": "2026-10-31T12:55:18Z"
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable explanation of error",
    "details": {
      "field": "Specific field error"
    }
  }
}
```

---

## 3. Implemented Phase 0 Endpoints

### Health Check
- **Endpoint**: `GET /api/v1/health`
- **Auth**: None
- **Response**: `200 OK`
  ```json
  {
    "database": "connected",
    "status": "healthy",
    "timestamp": "2026-10-31T07:11:10Z",
    "version": "0.1.0"
  }
  ```

### Customer Registration
- **Endpoint**: `POST /api/v1/auth/register`
- **Auth**: None
- **Request Body**:
  ```json
  {
    "email": "customer@company.com",
    "password": "StrongPassword123!",
    "companyName": "Acme Corp"
  }
  ```
- **Password Rules**:
  - Length strictly > 8 characters
  - At least 1 uppercase letter
  - At least 1 lowercase letter
  - At least 1 special character
- **Responses**:
  - `201 Created`: User created & session cookie issued.
  - `400 Bad Request`: `WEAK_PASSWORD` or `INVALID_EMAIL`.
  - `409 Conflict`: `DUPLICATE_EMAIL`.

### User Login
- **Endpoint**: `POST /api/v1/auth/login`
- **Auth**: None
- **Request Body**:
  ```json
  {
    "email": "admin@brightloop.local",
    "password": "AdminPassword123!"
  }
  ```
- **Responses**:
  - `200 OK`: Session cookie set.
  - `401 Unauthorized`: `INVALID_CREDENTIALS`.

### User Logout
- **Endpoint**: `POST /api/v1/auth/logout`
- **Auth**: Any active session
- **Response**: `200 OK` (cookie cleared).

### Current User Profile
- **Endpoint**: `GET /api/v1/auth/me`
- **Auth**: `RequireAuth` (any logged in user)
- **Responses**:
  - `200 OK`: Returns current user details.
  - `401 Unauthorized`: Missing or invalid session cookie.

### Admin Authorization Verification
- **Endpoint**: `GET /api/v1/admin/ping`
- **Auth**: `RequireAuth` + `RequireRole("ADMIN")`
- **Responses**:
  - `200 OK`: Returns confirmation for ADMIN role.
  - `403 Forbidden`: Authenticated user is not an ADMIN.

---

## 4. Future Feature Branch Endpoints (Phases 1-3)

### Member 1: Subscriptions (`feat/subscription-lifecycle`)
- `GET /api/v1/subscriptions`
- `POST /api/v1/subscriptions`
- `GET /api/v1/subscriptions/{id}`
- `POST /api/v1/subscriptions/{id}/pause`
- `POST /api/v1/subscriptions/{id}/resume`
- `POST /api/v1/subscriptions/{id}/cancel`
- `GET /api/v1/subscriptions/{id}/timeline`

### Member 2: Catalog & Plans (`feat/catalog-plans`)
- `GET /api/v1/products`
- `POST /api/v1/products`
- `GET /api/v1/plans`
- `POST /api/v1/plans`
- `GET /api/v1/discounts`
- `POST /api/v1/discounts`

### Member 3: Billing & Payments (`feat/billing-engine`)
- `POST /api/v1/billing/run`
- `GET /api/v1/billing/runs`
- `GET /api/v1/invoices`
- `GET /api/v1/invoices/{id}`
- `GET /api/v1/invoices/{id}/pdf`
- `POST /api/v1/invoices/{id}/payments`

### Member 4: Quotations & Customer Portal (`feat/quotations-portal`)
- `GET /api/v1/quotation-templates`
- `POST /api/v1/quotations`
- `POST /api/v1/quotations/{id}/convert`
- `GET /api/v1/reporting/dashboard`
