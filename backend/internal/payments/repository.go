package payments

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

func (r *pgRepository) RecordPayment(ctx context.Context, p *Payment) error {
	if p.PaymentDate.IsZero() {
		p.PaymentDate = time.Now().UTC()
	}

	query := `
		INSERT INTO payments (invoice_id, payment_number, amount_cents, payment_method, payment_date, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		p.InvoiceID, p.PaymentNumber, p.AmountCents, p.PaymentMethod, p.PaymentDate, p.Notes,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *pgRepository) ListPaymentsByInvoiceID(ctx context.Context, invoiceID string) ([]Payment, error) {
	query := `
		SELECT id, invoice_id, payment_number, amount_cents, payment_method, payment_date, notes, created_at
		FROM payments
		WHERE invoice_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]Payment, 0)
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.InvoiceID, &p.PaymentNumber, &p.AmountCents, &p.PaymentMethod, &p.PaymentDate, &p.Notes, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *pgRepository) ListAllPayments(ctx context.Context, customerID *string) ([]Payment, error) {
	query := `
		SELECT p.id, p.invoice_id, i.invoice_number, u.email,
		       p.payment_number, p.amount_cents, p.payment_method, p.payment_date, p.notes, p.created_at
		FROM payments p
		JOIN invoices i ON p.invoice_id = i.id
		JOIN users u ON i.customer_id = u.id
		WHERE ($1::uuid IS NULL OR i.customer_id = $1::uuid)
		ORDER BY p.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]Payment, 0)
	for rows.Next() {
		var p Payment
		if err := rows.Scan(
			&p.ID, &p.InvoiceID, &p.InvoiceNumber, &p.CustomerEmail,
			&p.PaymentNumber, &p.AmountCents, &p.PaymentMethod, &p.PaymentDate, &p.Notes, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *pgRepository) GetTotalPaidForInvoice(ctx context.Context, invoiceID string) (int64, error) {
	query := `SELECT COALESCE(SUM(amount_cents), 0) FROM payments WHERE invoice_id = $1`
	var total int64
	err := r.pool.QueryRow(ctx, query, invoiceID).Scan(&total)
	return total, err
}
