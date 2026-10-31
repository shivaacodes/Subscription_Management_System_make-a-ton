package catalog

import (
	"context"
	"time"
)

type Product struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	ProductType     string           `json:"productType"`
	SalesPriceCents int64            `json:"salesPriceCents"`
	CostPriceCents  int64            `json:"costPriceCents"`
	Notes           *string          `json:"notes,omitempty"`
	IsActive        bool             `json:"isActive"`
	Variants        []ProductVariant `json:"variants,omitempty"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type ProductVariant struct {
	ID                   string    `json:"id"`
	ProductID            string    `json:"productId"`
	Name                 string    `json:"name"`
	PriceAdjustmentCents int64     `json:"priceAdjustmentCents"`
	IsActive             bool      `json:"isActive"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type TaxRule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	RateBps   int       `json:"rateBps"` // e.g. 1800 for 18.00%
	Country   string    `json:"country"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type DiscountRule struct {
	ID                string     `json:"id"`
	Code              string     `json:"code"`
	DiscountType      string     `json:"discountType"` // FIXED_AMOUNT or PERCENTAGE
	ValueCents        int64      `json:"valueCents"`   // cents or percentage basis points
	MinPurchaseCents  int64      `json:"minPurchaseCents"`
	MinQuantity       int        `json:"minQuantity"`
	MaxUsageCount     *int       `json:"maxUsageCount,omitempty"`
	CurrentUsageCount int        `json:"currentUsageCount"`
	AppliesTo         string     `json:"appliesTo"` // PRODUCT or SUBSCRIPTION
	ProductID         *string    `json:"productId,omitempty"`
	ValidFrom         time.Time  `json:"validFrom"`
	ValidTo           *time.Time `json:"validTo,omitempty"`
	IsActive          bool       `json:"isActive"`
	CreatedAt         time.Time  `json:"createdAt"`
}

type Repository interface {
	CreateProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, activeOnly bool) ([]Product, error)

	CreateVariant(ctx context.Context, v *ProductVariant) error
	UpdateVariant(ctx context.Context, v *ProductVariant) error
	ListVariantsByProductID(ctx context.Context, productID string) ([]ProductVariant, error)

	CreateTaxRule(ctx context.Context, t *TaxRule) error
	ListTaxRules(ctx context.Context) ([]TaxRule, error)

	CreateDiscountRule(ctx context.Context, d *DiscountRule) error
	ListDiscountRules(ctx context.Context) ([]DiscountRule, error)
	GetDiscountRuleByCode(ctx context.Context, code string) (*DiscountRule, error)
}
