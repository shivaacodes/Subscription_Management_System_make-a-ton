package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/subflow/backend/internal/auth"
	"github.com/subflow/backend/internal/config"
	"github.com/subflow/backend/internal/platform"
	"github.com/subflow/backend/internal/users"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("running development database seed")
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := platform.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := platform.RunMigrations(ctx, pool, "migrations"); err != nil {
		slog.Error("migration failed during seed", "error", err)
		os.Exit(1)
	}

	userRepo := users.NewRepository(pool)

	// Seed Users
	seedUsers := []struct {
		Email       string
		Password    string
		Role        users.Role
		CompanyName string
	}{
		{
			Email:       "admin@brightloop.local",
			Password:    "AdminPassword123!",
			Role:        users.RoleAdmin,
			CompanyName: "Brightloop Software HQ",
		},
		{
			Email:       "staff@brightloop.local",
			Password:    "StaffPassword123!",
			Role:        users.RoleStaff,
			CompanyName: "Brightloop Operations",
		},
		{
			Email:       "customer@brightloop.local",
			Password:    "CustomerPassword123!",
			Role:        users.RoleCustomer,
			CompanyName: "Acme Innovations Ltd",
		},
	}

	for _, su := range seedUsers {
		existing, err := userRepo.GetByEmail(ctx, su.Email)
		if err == nil && existing != nil {
			slog.Info("user already exists, skipping", "email", su.Email)
			continue
		}

		hash, err := auth.HashPassword(su.Password)
		if err != nil {
			slog.Error("failed to hash password", "error", err)
			continue
		}

		u, _, err := userRepo.CreateUserWithProfile(ctx, su.Email, hash, su.Role, &su.CompanyName)
		if err != nil {
			slog.Error("failed to seed user", "email", su.Email, "error", err)
			continue
		}
		slog.Info("seeded user successfully", "email", u.Email, "role", u.Role)
	}

	// Seed Sample Product and Plan for Brightloop reference
	var productID string
	err = pool.QueryRow(ctx, `
		INSERT INTO products (name, product_type, sales_price_cents, cost_price_cents, notes, is_active)
		VALUES ('Brightloop Core Subscription', 'SUBSCRIPTION_SERVICE', 499900, 120000, 'Core platform subscription', true)
		ON CONFLICT DO NOTHING
		RETURNING id;
	`).Scan(&productID)
	if err == nil {
		slog.Info("seeded base product", "id", productID)

		// Seed Odoo Variant (+₹560 = 56000 cents)
		_, _ = pool.Exec(ctx, `
			INSERT INTO product_variants (product_id, name, price_adjustment_cents, is_active)
			VALUES ($1, 'Odoo Brand Connector', 56000, true);
		`, productID)
		slog.Info("seeded product variant (+₹560)")
	}

	// Seed Standard Plan
	_, _ = pool.Exec(ctx, `
		INSERT INTO plans (name, frequency, billing_price_cents, min_quantity, available_from, auto_close, allow_early_close, allow_pause, allow_renewal, is_active)
		VALUES ('Standard Monthly Plan', 'MONTHLY', 499900, 1, '2026-01-01', false, true, true, true, true)
		ON CONFLICT DO NOTHING;
	`)
	slog.Info("seeded monthly plan")

	fmt.Println("\nSeed completed successfully!")
	fmt.Println("Admin credentials:    admin@brightloop.local    / AdminPassword123!")
	fmt.Println("Staff credentials:    staff@brightloop.local    / StaffPassword123!")
	fmt.Println("Customer credentials: customer@brightloop.local / CustomerPassword123!")
}
