package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/subflow/backend/internal/auth"
	"github.com/subflow/backend/internal/billing"
	"github.com/subflow/backend/internal/catalog"
	"github.com/subflow/backend/internal/config"
	"github.com/subflow/backend/internal/payments"
	"github.com/subflow/backend/internal/plans"
	"github.com/subflow/backend/internal/platform"
	"github.com/subflow/backend/internal/quotations"
	"github.com/subflow/backend/internal/reporting"
	"github.com/subflow/backend/internal/subscriptions"
	"github.com/subflow/backend/internal/users"
)

func main() {
	// 1. Structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting Veyra (Subflow) API server")

	// 2. Load Configuration
	cfg := config.Load()

	// 3. Database Connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := platform.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 4. Run Migrations
	migrationsDir := findMigrationsDir()
	if err := platform.RunMigrations(ctx, pool, migrationsDir); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	// 5. Initialize Repositories and Services
	userRepo := users.NewRepository(pool)
	sessionRepo := auth.NewSessionRepository(pool)
	authService := auth.NewService(userRepo, sessionRepo)
	authHandler := auth.NewHandler(authService, cfg)

	catalogRepo := catalog.NewRepository(pool)
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogService)

	planRepo := plans.NewRepository(pool)
	planService := plans.NewService(planRepo)
	planHandler := plans.NewHandler(planService)

	subRepo := subscriptions.NewRepository(pool)
	subService := subscriptions.NewService(subRepo, planService)
	subHandler := subscriptions.NewHandler(subService)

	quoteRepo := quotations.NewRepository(pool)
	quoteService := quotations.NewService(quoteRepo, subService)
	quoteHandler := quotations.NewHandler(quoteService)

	billingRepo := billing.NewRepository(pool)
	billingService := billing.NewService(billingRepo, subRepo, planService)
	billingHandler := billing.NewHandler(billingService)

	paymentRepo := payments.NewRepository(pool)
	paymentService := payments.NewService(paymentRepo, billingService)
	paymentHandler := payments.NewHandler(paymentService)

	reportingRepo := reporting.NewRepository(pool)
	reportingService := reporting.NewService(reportingRepo)
	reportingHandler := reporting.NewHandler(reportingService)

	// 6. Router Setup
	r := chi.NewRouter()

	r.Use(platform.Recoverer)
	r.Use(platform.RequestLogger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(platform.SessionAuthMiddleware(authService))

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			dbStatus := "connected"
			if err := pool.Ping(r.Context()); err != nil {
				dbStatus = "unreachable"
			}
			platform.WriteJSON(w, http.StatusOK, map[string]string{
				"status":    "healthy",
				"database":  dbStatus,
				"version":   "0.1.0",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})

		// Authentication
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)

			r.Group(func(r chi.Router) {
				r.Use(platform.RequireAuth)
				r.Get("/me", authHandler.Me)
			})
		})

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(platform.RequireAuth)

			// Catalog & Products
			r.Route("/catalog", func(r chi.Router) {
				r.Get("/products", catalogHandler.ListProducts)
				r.Get("/products/{id}", catalogHandler.GetProduct)

				r.Group(func(r chi.Router) {
					r.Use(platform.RequireRole(users.RoleAdmin))
					r.Post("/products", catalogHandler.CreateProduct)
					r.Put("/products/{id}", catalogHandler.UpdateProduct)
					r.Post("/products/{productId}/variants", catalogHandler.CreateVariant)
					r.Put("/variants/{id}", catalogHandler.UpdateVariant)
				})
			})

			// Taxes & Discounts (Admin only)
			r.Route("/taxes", func(r chi.Router) {
				r.Get("/", catalogHandler.ListTaxes)
				r.Group(func(r chi.Router) {
					r.Use(platform.RequireRole(users.RoleAdmin))
					r.Post("/", catalogHandler.CreateTax)
				})
			})

			r.Route("/discounts", func(r chi.Router) {
				r.Get("/", catalogHandler.ListDiscounts)
				r.Group(func(r chi.Router) {
					r.Use(platform.RequireRole(users.RoleAdmin))
					r.Post("/", catalogHandler.CreateDiscount)
				})
			})

			// Subscription Plans
			r.Route("/plans", func(r chi.Router) {
				r.Get("/", planHandler.ListPlans)
				r.Get("/{id}", planHandler.GetPlan)

				r.Group(func(r chi.Router) {
					r.Use(platform.RequireRole(users.RoleAdmin))
					r.Post("/", planHandler.CreatePlan)
					r.Post("/{planId}/products", planHandler.LinkProduct)
				})
			})

			// Quotations
			r.Route("/quotations", func(r chi.Router) {
				r.Get("/", quoteHandler.ListQuotations)
				r.Post("/", quoteHandler.CreateQuotation)
				r.Get("/{id}", quoteHandler.GetQuotation)
				r.Post("/{id}/convert", quoteHandler.ConvertQuotation)
			})
			r.Get("/quotation-templates", quoteHandler.ListTemplates)

			// Subscriptions
			r.Route("/subscriptions", func(r chi.Router) {
				r.Get("/", subHandler.ListSubscriptions)
				r.Post("/", subHandler.CreateSubscription)
				r.Get("/{id}", subHandler.GetSubscription)
				r.Post("/{id}/activate", subHandler.ActivateSubscription)
				r.Post("/{id}/pause", subHandler.PauseSubscription)
				r.Post("/{id}/resume", subHandler.ResumeSubscription)
				r.Post("/{id}/cancel", subHandler.CancelSubscription)
				r.Post("/{id}/addons", subHandler.ToggleAddon)
				r.Get("/{id}/timeline", subHandler.GetTimeline)
			})

			// Billing & Invoices
			r.Route("/billing", func(r chi.Router) {
				r.Post("/preview", billingHandler.PreviewCalculation)
				r.Get("/reliability", billingHandler.GetReliabilityCenter)
				r.Post("/run", billingHandler.TriggerBillingRun)
			})

			r.Route("/invoices", func(r chi.Router) {
				r.Get("/", billingHandler.ListInvoices)
				r.Get("/{id}", billingHandler.GetInvoice)
				r.Post("/{id}/confirm", billingHandler.ConfirmInvoice)
				r.Post("/{invoiceId}/payments", paymentHandler.RecordPayment)
			})

			// Payments
			r.Route("/payments", func(r chi.Router) {
				r.Get("/", paymentHandler.ListPayments)
			})

			// Reporting & Analytics
			r.Route("/reporting", func(r chi.Router) {
				r.Use(platform.RequireRole(users.RoleAdmin, users.RoleStaff))
				r.Get("/dashboard", reportingHandler.GetDashboardMetrics)
			})

			// Admin Verification
			r.Group(func(r chi.Router) {
				r.Use(platform.RequireRole(users.RoleAdmin))
				r.Get("/admin/ping", func(w http.ResponseWriter, r *http.Request) {
					user := platform.GetUserFromContext(r.Context())
					platform.WriteJSON(w, http.StatusOK, map[string]string{
						"message":   fmt.Sprintf("Welcome Admin %s, role authorization successfully verified", user.Email),
						"timestamp": time.Now().UTC().Format(time.RFC3339),
					})
				})
			})
		})
	})

	// 7. HTTP Server with graceful shutdown
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server listening", "port", cfg.Port, "env", cfg.Env)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		slog.Info("shutdown signal received", "signal", sig.String())
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(ctxShutdown); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			_ = server.Close()
		}
		slog.Info("server stopped gracefully")
	}
}

func findMigrationsDir() string {
	candidates := []string{
		"migrations",
		"backend/migrations",
		"../migrations",
		filepath.Join(".", "migrations"),
	}
	for _, c := range candidates {
		if stat, err := os.Stat(c); err == nil && stat.IsDir() {
			return c
		}
	}
	return "migrations"
}
