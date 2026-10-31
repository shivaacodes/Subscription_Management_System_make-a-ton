package billing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvoiceNotFound = errors.New("invoice not found")

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) CreateInvoice(ctx context.Context, inv *Invoice, lines []InvoiceLine) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO invoices (
			invoice_number, subscription_id, customer_id, billing_run_id,
			status, billing_period_start, billing_period_end, currency,
			subtotal_cents, tax_cents, discount_cents, total_cents, due_date,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		inv.InvoiceNumber, inv.SubscriptionID, inv.CustomerID, inv.BillingRunID,
		inv.Status, inv.BillingPeriodStart, inv.BillingPeriodEnd, inv.Currency,
		inv.SubtotalCents, inv.TaxCents, inv.DiscountCents, inv.TotalCents, inv.DueDate,
	).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return err
	}

	lineQuery := `
		INSERT INTO invoice_lines (
			invoice_id, product_id, variant_id, description, quantity,
			unit_price_cents, tax_cents, discount_cents, total_cents, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`
	for i := range lines {
		lines[i].InvoiceID = inv.ID
		err = tx.QueryRow(ctx, lineQuery,
			inv.ID, lines[i].ProductID, lines[i].VariantID, lines[i].Description,
			lines[i].Quantity, lines[i].UnitPriceCents, lines[i].TaxCents, lines[i].DiscountCents,
			lines[i].TotalCents,
		).Scan(&lines[i].ID, &lines[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("failed inserting invoice line %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) GetInvoiceByID(ctx context.Context, id string) (*Invoice, []InvoiceLine, error) {
	query := `
		SELECT i.id, i.invoice_number, i.subscription_id, s.subscription_number,
		       i.customer_id, u.email, i.billing_run_id, i.status,
		       i.billing_period_start, i.billing_period_end, i.currency,
		       i.subtotal_cents, i.tax_cents, i.discount_cents, i.total_cents,
		       COALESCE((SELECT SUM(amount_cents) FROM payments WHERE invoice_id = i.id), 0) as paid_amount,
		       i.due_date, i.paid_at, i.created_at, i.updated_at
		FROM invoices i
		JOIN subscriptions s ON i.subscription_id = s.id
		JOIN users u ON i.customer_id = u.id
		WHERE i.id = $1
	`
	var inv Invoice
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.InvoiceNumber, &inv.SubscriptionID, &inv.SubscriptionNumber,
		&inv.CustomerID, &inv.CustomerEmail, &inv.BillingRunID, &inv.Status,
		&inv.BillingPeriodStart, &inv.BillingPeriodEnd, &inv.Currency,
		&inv.SubtotalCents, &inv.TaxCents, &inv.DiscountCents, &inv.TotalCents,
		&inv.PaidAmountCents, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrInvoiceNotFound
		}
		return nil, nil, err
	}
	inv.OutstandingBalanceCents = inv.TotalCents - inv.PaidAmountCents

	lineQuery := `
		SELECT il.id, il.invoice_id, il.product_id, pr.name, il.variant_id, COALESCE(pv.name, ''),
		       il.description, il.quantity, il.unit_price_cents, il.tax_cents, il.discount_cents, il.total_cents, il.created_at
		FROM invoice_lines il
		JOIN products pr ON il.product_id = pr.id
		LEFT JOIN product_variants pv ON il.variant_id = pv.id
		WHERE il.invoice_id = $1
	`
	rows, err := r.pool.Query(ctx, lineQuery, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var lines []InvoiceLine
	for rows.Next() {
		var l InvoiceLine
		var variantName string
		if err := rows.Scan(
			&l.ID, &l.InvoiceID, &l.ProductID, &l.ProductName, &l.VariantID, &variantName,
			&l.Description, &l.Quantity, &l.UnitPriceCents, &l.TaxCents, &l.DiscountCents, &l.TotalCents, &l.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		if variantName != "" {
			l.VariantName = &variantName
		}
		lines = append(lines, l)
	}

	inv.Lines = lines
	return &inv, lines, nil
}

func (r *pgRepository) ListInvoices(ctx context.Context, customerID *string, status *string) ([]Invoice, error) {
	query := `
		SELECT i.id, i.invoice_number, i.subscription_id, s.subscription_number,
		       i.customer_id, u.email, i.billing_run_id, i.status,
		       i.billing_period_start, i.billing_period_end, i.currency,
		       i.subtotal_cents, i.tax_cents, i.discount_cents, i.total_cents,
		       COALESCE((SELECT SUM(amount_cents) FROM payments WHERE invoice_id = i.id), 0) as paid_amount,
		       i.due_date, i.paid_at, i.created_at, i.updated_at
		FROM invoices i
		JOIN subscriptions s ON i.subscription_id = s.id
		JOIN users u ON i.customer_id = u.id
		WHERE ($1::uuid IS NULL OR i.customer_id = $1::uuid)
		  AND ($2::varchar IS NULL OR i.status = $2::varchar)
		ORDER BY i.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, customerID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Invoice
	for rows.Next() {
		var inv Invoice
		if err := rows.Scan(
			&inv.ID, &inv.InvoiceNumber, &inv.SubscriptionID, &inv.SubscriptionNumber,
			&inv.CustomerID, &inv.CustomerEmail, &inv.BillingRunID, &inv.Status,
			&inv.BillingPeriodStart, &inv.BillingPeriodEnd, &inv.Currency,
			&inv.SubtotalCents, &inv.TaxCents, &inv.DiscountCents, &inv.TotalCents,
			&inv.PaidAmountCents, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inv.OutstandingBalanceCents = inv.TotalCents - inv.PaidAmountCents
		list = append(list, inv)
	}
	return list, nil
}

func (r *pgRepository) UpdateInvoiceStatus(ctx context.Context, id string, status InvoiceStatus) error {
	query := `UPDATE invoices SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *pgRepository) InvoiceExistsForPeriod(ctx context.Context, subscriptionID string, start, end time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM invoices
			WHERE subscription_id = $1 AND billing_period_start = $2 AND billing_period_end = $3
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, query, subscriptionID, start, end).Scan(&exists)
	return exists, err
}

func (r *pgRepository) CreateBillingRun(ctx context.Context, run *BillingRun) error {
	query := `
		INSERT INTO billing_runs (run_date, status, subscriptions_evaluated, invoices_generated, failures_count, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, run.RunDate, run.Status, run.SubscriptionsEvaluated, run.InvoicesGenerated, run.FailuresCount).
		Scan(&run.ID, &run.CreatedAt)
}

func (r *pgRepository) CompleteBillingRun(ctx context.Context, id string, evaluated, generated, failures int, errorSummary *string) error {
	query := `
		UPDATE billing_runs
		SET status = 'COMPLETED', subscriptions_evaluated = $1, invoices_generated = $2, failures_count = $3,
		    error_summary = $4, completed_at = NOW()
		WHERE id = $5
	`
	_, err := r.pool.Exec(ctx, query, evaluated, generated, failures, errorSummary, id)
	return err
}

func (r *pgRepository) RecordBillingAttempt(ctx context.Context, attempt *BillingAttempt) error {
	query := `
		INSERT INTO billing_attempts (billing_run_id, subscription_id, status, error_message, attempt_number, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, attempt.BillingRunID, attempt.SubscriptionID, attempt.Status, attempt.ErrorMessage, attempt.AttemptNumber).
		Scan(&attempt.ID, &attempt.CreatedAt)
}

func (r *pgRepository) ListRecentBillingRuns(ctx context.Context, limit int) ([]BillingRun, error) {
	query := `
		SELECT id, run_date, status, subscriptions_evaluated, invoices_generated, failures_count, error_summary, created_at, completed_at
		FROM billing_runs
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []BillingRun
	for rows.Next() {
		var run BillingRun
		if err := rows.Scan(
			&run.ID, &run.RunDate, &run.Status, &run.SubscriptionsEvaluated, &run.InvoicesGenerated,
			&run.FailuresCount, &run.ErrorSummary, &run.CreatedAt, &run.CompletedAt,
		); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func (r *pgRepository) ListRecentAttempts(ctx context.Context, limit int) ([]BillingAttempt, error) {
	query := `
		SELECT ba.id, ba.billing_run_id, ba.subscription_id, s.subscription_number,
		       ba.status, ba.error_message, ba.attempt_number, ba.created_at
		FROM billing_attempts ba
		JOIN subscriptions s ON ba.subscription_id = s.id
		ORDER BY ba.created_at DESC
		LIMIT $1
	`
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []BillingAttempt
	for rows.Next() {
		var a BillingAttempt
		if err := rows.Scan(&a.ID, &a.BillingRunID, &a.SubscriptionID, &a.SubscriptionNumber, &a.Status, &a.ErrorMessage, &a.AttemptNumber, &a.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}
