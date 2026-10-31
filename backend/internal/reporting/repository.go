package reporting

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) GetDashboardMetrics(ctx context.Context, asOf time.Time) (*DashboardMetrics, error) {
	var metrics DashboardMetrics

	// 1. Active Subscriptions count
	_ = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM subscriptions WHERE status = 'ACTIVE'
	`).Scan(&metrics.ActiveSubscriptionsCount)

	// 2. MRR (Monthly Recurring Revenue)
	// Normalized calculation: Daily * 30, Weekly * 4, Monthly * 1, Yearly / 12
	_ = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(
			CASE
				WHEN p.frequency = 'DAILY' THEN p.billing_price_cents * 30
				WHEN p.frequency = 'WEEKLY' THEN p.billing_price_cents * 4
				WHEN p.frequency = 'MONTHLY' THEN p.billing_price_cents
				WHEN p.frequency = 'YEARLY' THEN p.billing_price_cents / 12
				ELSE p.billing_price_cents
			END
		), 0)
		FROM subscriptions s
		JOIN plans p ON s.plan_id = p.id
		WHERE s.status = 'ACTIVE'
	`).Scan(&metrics.MRRCents)

	// 3. Invoiced Revenue
	_ = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_cents), 0)
		FROM invoices
		WHERE status != 'CANCELLED'
	`).Scan(&metrics.InvoicedRevenueCents)

	// 4. Collected Payments
	_ = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_cents), 0)
		FROM payments
	`).Scan(&metrics.CollectedPaymentsCents)

	// 5. Outstanding Receivables
	metrics.OutstandingReceivables = metrics.InvoicedRevenueCents - metrics.CollectedPaymentsCents
	if metrics.OutstandingReceivables < 0 {
		metrics.OutstandingReceivables = 0
	}

	// 6. Overdue Invoices Count
	_ = r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM invoices
		WHERE status != 'PAID' AND status != 'CANCELLED' AND due_date < $1
	`, asOf).Scan(&metrics.OverdueInvoicesCount)

	// 7. Monthly Trend (last 6 months)
	trendQuery := `
		SELECT TO_CHAR(months.m, 'Mon YYYY') as month_label,
		       COALESCE((SELECT SUM(total_cents) FROM invoices WHERE status != 'CANCELLED' AND DATE_TRUNC('month', created_at) = months.m), 0) as invoiced,
		       COALESCE((SELECT SUM(amount_cents) FROM payments WHERE DATE_TRUNC('month', payment_date) = months.m), 0) as collected
		FROM (
			SELECT DATE_TRUNC('month', NOW() - (s || ' months')::interval) as m
			FROM generate_series(5, 0, -1) as s
		) months
		ORDER BY months.m ASC
	`
	rows, err := r.pool.Query(ctx, trendQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var pt MonthlyTrendPoint
			if err := rows.Scan(&pt.Month, &pt.InvoicedRevenueCents, &pt.CollectedCents); err == nil {
				metrics.Trend = append(metrics.Trend, pt)
			}
		}
	}

	return &metrics, nil
}
