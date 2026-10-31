package quotations

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrQuotationNotFound = errors.New("quotation not found")

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) CreateQuotation(ctx context.Context, q *Quotation, lines []QuotationLine) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO quotations (
			quotation_number, customer_id, plan_id, status, valid_until,
			subtotal_cents, tax_cents, discount_cents, total_cents, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		q.QuotationNumber, q.CustomerID, q.PlanID, q.Status, q.ValidUntil,
		q.SubtotalCents, q.TaxCents, q.DiscountCents, q.TotalCents,
	).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return err
	}

	lineQuery := `
		INSERT INTO quotation_lines (
			quotation_id, product_id, variant_id, quantity, unit_price_cents, tax_rate_bps, total_cents
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	for i := range lines {
		lines[i].QuotationID = q.ID
		err = tx.QueryRow(ctx, lineQuery,
			q.ID, lines[i].ProductID, lines[i].VariantID, lines[i].Quantity,
			lines[i].UnitPriceCents, lines[i].TaxRateBps, lines[i].TotalCents,
		).Scan(&lines[i].ID)
		if err != nil {
			return fmt.Errorf("failed inserting line %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) GetQuotationByID(ctx context.Context, id string) (*Quotation, []QuotationLine, error) {
	query := `
		SELECT q.id, q.quotation_number, q.customer_id, u.email, q.plan_id, p.name,
		       q.status, q.valid_until, q.subtotal_cents, q.tax_cents, q.discount_cents, q.total_cents,
		       q.created_at, q.updated_at
		FROM quotations q
		JOIN users u ON q.customer_id = u.id
		JOIN plans p ON q.plan_id = p.id
		WHERE q.id = $1
	`
	var q Quotation
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.QuotationNumber, &q.CustomerID, &q.CustomerEmail, &q.PlanID, &q.PlanName,
		&q.Status, &q.ValidUntil, &q.SubtotalCents, &q.TaxCents, &q.DiscountCents, &q.TotalCents,
		&q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrQuotationNotFound
		}
		return nil, nil, err
	}

	lineQuery := `
		SELECT ql.id, ql.quotation_id, ql.product_id, pr.name, ql.variant_id, COALESCE(pv.name, ''),
		       ql.quantity, ql.unit_price_cents, ql.tax_rate_bps, ql.total_cents
		FROM quotation_lines ql
		JOIN products pr ON ql.product_id = pr.id
		LEFT JOIN product_variants pv ON ql.variant_id = pv.id
		WHERE ql.quotation_id = $1
	`
	rows, err := r.pool.Query(ctx, lineQuery, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var lines []QuotationLine
	for rows.Next() {
		var l QuotationLine
		var variantName string
		if err := rows.Scan(
			&l.ID, &l.QuotationID, &l.ProductID, &l.ProductName, &l.VariantID, &variantName,
			&l.Quantity, &l.UnitPriceCents, &l.TaxRateBps, &l.TotalCents,
		); err != nil {
			return nil, nil, err
		}
		if variantName != "" {
			l.VariantName = &variantName
		}
		lines = append(lines, l)
	}

	q.Lines = lines
	return &q, lines, nil
}

func (r *pgRepository) ListQuotations(ctx context.Context, customerID *string) ([]Quotation, error) {
	query := `
		SELECT q.id, q.quotation_number, q.customer_id, u.email, q.plan_id, p.name,
		       q.status, q.valid_until, q.subtotal_cents, q.tax_cents, q.discount_cents, q.total_cents,
		       q.created_at, q.updated_at
		FROM quotations q
		JOIN users u ON q.customer_id = u.id
		JOIN plans p ON q.plan_id = p.id
		WHERE ($1::uuid IS NULL OR q.customer_id = $1::uuid)
		ORDER BY q.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Quotation
	for rows.Next() {
		var q Quotation
		if err := rows.Scan(
			&q.ID, &q.QuotationNumber, &q.CustomerID, &q.CustomerEmail, &q.PlanID, &q.PlanName,
			&q.Status, &q.ValidUntil, &q.SubtotalCents, &q.TaxCents, &q.DiscountCents, &q.TotalCents,
			&q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, nil
}

func (r *pgRepository) UpdateStatus(ctx context.Context, id string, status QuotationStatus) error {
	query := `UPDATE quotations SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *pgRepository) CreateTemplate(ctx context.Context, t *QuotationTemplate, lines []QuotationTemplateLine) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO quotation_templates (name, description, validity_days, plan_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, t.Name, t.Description, t.ValidityDays, t.PlanID, t.IsActive).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return err
	}

	lineQuery := `
		INSERT INTO quotation_template_lines (template_id, product_id, variant_id, quantity, unit_price_cents)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, l := range lines {
		if _, err := tx.Exec(ctx, lineQuery, t.ID, l.ProductID, l.VariantID, l.Quantity, l.UnitPriceCents); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRepository) ListTemplates(ctx context.Context) ([]QuotationTemplate, error) {
	query := `
		SELECT id, name, description, validity_days, plan_id, is_active, created_at, updated_at
		FROM quotation_templates
		WHERE is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []QuotationTemplate
	for rows.Next() {
		var t QuotationTemplate
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.ValidityDays, &t.PlanID, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}
