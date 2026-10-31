# Local Development Guide

## Prerequisites
- **Go**: 1.23+ or 1.24+ (Standard toolchain)
- **Node.js**: 20+ (with npm)
- **PostgreSQL**: 16+ or 17 (or Docker Compose)

---

## Quickstart

### 1. Start PostgreSQL
Using Docker Compose:
```bash
docker-compose up -d postgres
```
*Or using local PostgreSQL service on `localhost:5432` with user `subflow` and database `subflow`.*

### 2. Configure Environment Variables
Copy example configs:
```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

### 3. Run Database Migrations and Seed Demo Data
```bash
cd backend
go run cmd/seed/main.go
```
This automatically applies migrations from `backend/migrations/` and provisions test credentials:
- **Admin**: `admin@brightloop.local` / `AdminPassword123!`
- **Staff**: `staff@brightloop.local` / `StaffPassword123!`
- **Customer**: `customer@brightloop.local` / `CustomerPassword123!`

### 4. Start the Backend API Server
```bash
cd backend
go run cmd/api/main.go
```
The server starts on `http://localhost:8080`.

Verify health check:
```bash
curl http://localhost:8080/api/v1/health
```

### 5. Start the Frontend Development Server
```bash
cd frontend
npm install
npm run dev
```
The application opens on `http://localhost:5173`.

---

## Running Automated Tests

### Backend Unit & Middleware Tests
```bash
cd backend
go test -v ./...
```

### Check Go Formatting & Vet
```bash
cd backend
test -z "$(gofmt -l .)" && echo "Formatted"
go vet ./...
```

### Frontend Typecheck & Build
```bash
cd frontend
npx tsc --noEmit
npm run build
```
