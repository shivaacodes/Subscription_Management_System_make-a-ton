package plans

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPlanNotFound = errors.New("plan not found")

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) CreatePlan(ctx context.Context, p *Plan) error {
	query := `
		INSERT INTO plans (
			name, frequency, billing_price_cents, min_quantity,
			available_from, available_to, auto_close, allow_early_close,
			allow_pause, allow_renewal, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		p.Name, p.Frequency, p.BillingPriceCents, p.MinQuantity,
		p.AvailableFrom, p.AvailableTo, p.AutoClose, p.AllowEarlyClose,
		p.AllowPause, p.AllowRenewal, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *pgRepository) GetPlanByID(ctx context.Context, id string) (*Plan, error) {
	query := `
		SELECT id, name, frequency, billing_price_cents, min_quantity,
		       available_from, available_to, auto_close, allow_early_close,
		       allow_pause, allow_renewal, is_active, created_at, updated_at
		FROM plans
		WHERE id = $1
	`
	var p Plan
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Frequency, &p.BillingPriceCents, &p.MinQuantity,
		&p.AvailableFrom, &p.AvailableTo, &p.AutoClose, &p.AllowEarlyClose,
		&p.AllowPause, &p.AllowRenewal, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgRepository) ListActivePlans(ctx context.Context) ([]Plan, error) {
	query := `
		SELECT id, name, frequency, billing_price_cents, min_quantity,
		       available_from, available_to, auto_close, allow_early_close,
		       allow_pause, allow_renewal, is_active, created_at, updated_at
		FROM plans
		WHERE is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]Plan, 0)
	for rows.Next() {
		var p Plan
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Frequency, &p.BillingPriceCents, &p.MinQuantity,
			&p.AvailableFrom, &p.AvailableTo, &p.AutoClose, &p.AllowEarlyClose,
			&p.AllowPause, &p.AllowRenewal, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *pgRepository) AssociateProduct(ctx context.Context, pp *PlanProduct) error {
	query := `
		INSERT INTO plan_products (plan_id, product_id, is_addon, is_required, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (plan_id, product_id) DO UPDATE SET is_addon = $3, is_required = $4
		RETURNING id
	`
	return r.pool.QueryRow(ctx, query, pp.PlanID, pp.ProductID, pp.IsAddon, pp.IsRequired).Scan(&pp.ID)
}

func (r *pgRepository) GetPlanProducts(ctx context.Context, planID string) ([]PlanProduct, error) {
	query := `
		SELECT id, plan_id, product_id, is_addon, is_required
		FROM plan_products
		WHERE plan_id = $1
	`
	rows, err := r.pool.Query(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]PlanProduct, 0)
	for rows.Next() {
		var pp PlanProduct
		if err := rows.Scan(&pp.ID, &pp.PlanID, &pp.ProductID, &pp.IsAddon, &pp.IsRequired); err != nil {
			return nil, err
		}
		products = append(products, pp)
	}
	return products, nil
}
