package catalog

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrVariantNotFound  = errors.New("variant not found")
	ErrDiscountNotFound = errors.New("discount rule not found")
)

type pgRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{pool: pool}
}

func (r *pgRepository) CreateProduct(ctx context.Context, p *Product) error {
	query := `
		INSERT INTO products (name, product_type, sales_price_cents, cost_price_cents, notes, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, p.Name, p.ProductType, p.SalesPriceCents, p.CostPriceCents, p.Notes, p.IsActive).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *pgRepository) UpdateProduct(ctx context.Context, p *Product) error {
	query := `
		UPDATE products
		SET name = $1, product_type = $2, sales_price_cents = $3, cost_price_cents = $4, notes = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query, p.Name, p.ProductType, p.SalesPriceCents, p.CostPriceCents, p.Notes, p.IsActive, p.ID).
		Scan(&p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrProductNotFound
	}
	return err
}

func (r *pgRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT id, name, product_type, sales_price_cents, cost_price_cents, notes, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	var p Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.ProductType, &p.SalesPriceCents, &p.CostPriceCents, &p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	variants, err := r.ListVariantsByProductID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Variants = variants

	return &p, nil
}

func (r *pgRepository) ListProducts(ctx context.Context, activeOnly bool) ([]Product, error) {
	query := `
		SELECT id, name, product_type, sales_price_cents, cost_price_cents, notes, is_active, created_at, updated_at
		FROM products
		WHERE ($1 = false OR is_active = true)
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, activeOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.ProductType, &p.SalesPriceCents, &p.CostPriceCents, &p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	for i := range products {
		variants, err := r.ListVariantsByProductID(ctx, products[i].ID)
		if err == nil {
			products[i].Variants = variants
		}
	}

	return products, nil
}

func (r *pgRepository) CreateVariant(ctx context.Context, v *ProductVariant) error {
	query := `
		INSERT INTO product_variants (product_id, name, price_adjustment_cents, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, v.ProductID, v.Name, v.PriceAdjustmentCents, v.IsActive).
		Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

func (r *pgRepository) UpdateVariant(ctx context.Context, v *ProductVariant) error {
	query := `
		UPDATE product_variants
		SET name = $1, price_adjustment_cents = $2, is_active = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query, v.Name, v.PriceAdjustmentCents, v.IsActive, v.ID).Scan(&v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrVariantNotFound
	}
	return err
}

func (r *pgRepository) ListVariantsByProductID(ctx context.Context, productID string) ([]ProductVariant, error) {
	query := `
		SELECT id, product_id, name, price_adjustment_cents, is_active, created_at, updated_at
		FROM product_variants
		WHERE product_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := make([]ProductVariant, 0)
	for rows.Next() {
		var v ProductVariant
		if err := rows.Scan(&v.ID, &v.ProductID, &v.Name, &v.PriceAdjustmentCents, &v.IsActive, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}
	return variants, nil
}

func (r *pgRepository) CreateTaxRule(ctx context.Context, t *TaxRule) error {
	query := `
		INSERT INTO tax_rules (name, rate_bps, country, is_active, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, t.Name, t.RateBps, t.Country, t.IsActive).Scan(&t.ID, &t.CreatedAt)
}

func (r *pgRepository) ListTaxRules(ctx context.Context) ([]TaxRule, error) {
	query := `
		SELECT id, name, rate_bps, country, is_active, created_at
		FROM tax_rules
		WHERE is_active = true
		ORDER BY rate_bps DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]TaxRule, 0)
	for rows.Next() {
		var t TaxRule
		if err := rows.Scan(&t.ID, &t.Name, &t.RateBps, &t.Country, &t.IsActive, &t.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, t)
	}
	return rules, nil
}

func (r *pgRepository) CreateDiscountRule(ctx context.Context, d *DiscountRule) error {
	query := `
		INSERT INTO discount_rules (
			code, discount_type, value_cents, min_purchase_cents, min_quantity,
			max_usage_count, current_usage_count, applies_to, product_id,
			valid_from, valid_to, is_active, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8, $9, $10, $11, NOW())
		RETURNING id, current_usage_count, created_at
	`
	return r.pool.QueryRow(ctx, query,
		d.Code, d.DiscountType, d.ValueCents, d.MinPurchaseCents, d.MinQuantity,
		d.MaxUsageCount, d.AppliesTo, d.ProductID, d.ValidFrom, d.ValidTo, d.IsActive,
	).Scan(&d.ID, &d.CurrentUsageCount, &d.CreatedAt)
}

func (r *pgRepository) ListDiscountRules(ctx context.Context) ([]DiscountRule, error) {
	query := `
		SELECT id, code, discount_type, value_cents, min_purchase_cents, min_quantity,
		       max_usage_count, current_usage_count, applies_to, product_id, valid_from, valid_to, is_active, created_at
		FROM discount_rules
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]DiscountRule, 0)
	for rows.Next() {
		var d DiscountRule
		if err := rows.Scan(
			&d.ID, &d.Code, &d.DiscountType, &d.ValueCents, &d.MinPurchaseCents, &d.MinQuantity,
			&d.MaxUsageCount, &d.CurrentUsageCount, &d.AppliesTo, &d.ProductID, &d.ValidFrom, &d.ValidTo, &d.IsActive, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		rules = append(rules, d)
	}
	return rules, nil
}

func (r *pgRepository) GetDiscountRuleByCode(ctx context.Context, code string) (*DiscountRule, error) {
	query := `
		SELECT id, code, discount_type, value_cents, min_purchase_cents, min_quantity,
		       max_usage_count, current_usage_count, applies_to, product_id, valid_from, valid_to, is_active, created_at
		FROM discount_rules
		WHERE UPPER(code) = UPPER($1) AND is_active = true
	`
	var d DiscountRule
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&d.ID, &d.Code, &d.DiscountType, &d.ValueCents, &d.MinPurchaseCents, &d.MinQuantity,
		&d.MaxUsageCount, &d.CurrentUsageCount, &d.AppliesTo, &d.ProductID, &d.ValidFrom, &d.ValidTo, &d.IsActive, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDiscountNotFound
		}
		return nil, fmt.Errorf("failed querying discount rule: %w", err)
	}
	return &d, nil
}
