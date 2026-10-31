package platform

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/subflow/backend/internal/users"
)

func TestRequireAuthMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := RequireAuth(nextHandler)

	// Case 1: unauthenticated
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
	}

	// Case 2: authenticated user in context
	user := &users.User{
		ID:        "11111111-1111-1111-1111-111111111111",
		Email:     "user@example.com",
		Role:      users.RoleCustomer,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ctx := context.WithValue(context.Background(), UserContextKey, user)
	reqWithAuth := httptest.NewRequest("GET", "/test", nil).WithContext(ctx)
	recWithAuth := httptest.NewRecorder()
	handler.ServeHTTP(recWithAuth, reqWithAuth)

	if recWithAuth.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recWithAuth.Code)
	}
}

func TestRequireRoleMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("admin_ok"))
	})

	adminOnly := RequireRole(users.RoleAdmin)(nextHandler)

	customerUser := &users.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Email: "customer@example.com",
		Role:  users.RoleCustomer,
	}

	adminUser := &users.User{
		ID:    "22222222-2222-2222-2222-222222222222",
		Email: "admin@example.com",
		Role:  users.RoleAdmin,
	}

	// Case 1: customer tries accessing admin route
	ctxCustomer := context.WithValue(context.Background(), UserContextKey, customerUser)
	reqCustomer := httptest.NewRequest("GET", "/admin/test", nil).WithContext(ctxCustomer)
	recCustomer := httptest.NewRecorder()
	adminOnly.ServeHTTP(recCustomer, reqCustomer)

	if recCustomer.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for customer on admin route, got %d", recCustomer.Code)
	}

	// Case 2: admin access
	ctxAdmin := context.WithValue(context.Background(), UserContextKey, adminUser)
	reqAdmin := httptest.NewRequest("GET", "/admin/test", nil).WithContext(ctxAdmin)
	recAdmin := httptest.NewRecorder()
	adminOnly.ServeHTTP(recAdmin, reqAdmin)

	if recAdmin.Code != http.StatusOK {
		t.Errorf("expected 200 OK for admin on admin route, got %d", recAdmin.Code)
	}
}
