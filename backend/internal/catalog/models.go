package catalog

import (
	"context"
	"time"
)

type Product struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	ProductType     string    `json:"productType"`
	SalesPriceCents int64     `json:"salesPriceCents"`
	CostPriceCents  int64     `json:"costPriceCents"`
	Notes           *string   `json:"notes,omitempty"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
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

type Repository interface {
	CreateProduct(ctx context.Context, p *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListActiveProducts(ctx context.Context) ([]Product, error)
	CreateVariant(ctx context.Context, v *ProductVariant) error
	ListVariantsByProductID(ctx context.Context, productID string) ([]ProductVariant, error)
}
