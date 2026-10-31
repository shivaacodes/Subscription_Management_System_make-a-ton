package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidPrice = errors.New("sales price must be non-negative")
	ErrEmptyName    = errors.New("name cannot be empty")
)

type Service interface {
	CreateProduct(ctx context.Context, p *Product) error
	UpdateProduct(ctx context.Context, p *Product) error
	GetProduct(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, activeOnly bool) ([]Product, error)

	CreateVariant(ctx context.Context, v *ProductVariant) error
	UpdateVariant(ctx context.Context, v *ProductVariant) error

	CreateTaxRule(ctx context.Context, t *TaxRule) error
	ListTaxRules(ctx context.Context) ([]TaxRule, error)

	CreateDiscountRule(ctx context.Context, d *DiscountRule) error
	ListDiscountRules(ctx context.Context) ([]DiscountRule, error)
}

type catalogService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &catalogService{repo: repo}
}

func (s *catalogService) CreateProduct(ctx context.Context, p *Product) error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.SalesPriceCents < 0 || p.CostPriceCents < 0 {
		return ErrInvalidPrice
	}
	if p.ProductType == "" {
		p.ProductType = "SUBSCRIPTION_SERVICE"
	}
	return s.repo.CreateProduct(ctx, p)
}

func (s *catalogService) UpdateProduct(ctx context.Context, p *Product) error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.SalesPriceCents < 0 || p.CostPriceCents < 0 {
		return ErrInvalidPrice
	}
	return s.repo.UpdateProduct(ctx, p)
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *catalogService) ListProducts(ctx context.Context, activeOnly bool) ([]Product, error) {
	return s.repo.ListProducts(ctx, activeOnly)
}

func (s *catalogService) CreateVariant(ctx context.Context, v *ProductVariant) error {
	v.Name = strings.TrimSpace(v.Name)
	if v.Name == "" {
		return ErrEmptyName
	}
	return s.repo.CreateVariant(ctx, v)
}

func (s *catalogService) UpdateVariant(ctx context.Context, v *ProductVariant) error {
	v.Name = strings.TrimSpace(v.Name)
	if v.Name == "" {
		return ErrEmptyName
	}
	return s.repo.UpdateVariant(ctx, v)
}

func (s *catalogService) CreateTaxRule(ctx context.Context, t *TaxRule) error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return ErrEmptyName
	}
	if t.RateBps < 0 {
		return fmt.Errorf("tax rate cannot be negative")
	}
	if t.Country == "" {
		t.Country = "IN"
	}
	return s.repo.CreateTaxRule(ctx, t)
}

func (s *catalogService) ListTaxRules(ctx context.Context) ([]TaxRule, error) {
	return s.repo.ListTaxRules(ctx)
}

func (s *catalogService) CreateDiscountRule(ctx context.Context, d *DiscountRule) error {
	d.Code = strings.ToUpper(strings.TrimSpace(d.Code))
	if d.Code == "" {
		return fmt.Errorf("discount code cannot be empty")
	}
	if d.ValueCents <= 0 {
		return fmt.Errorf("discount value must be greater than zero")
	}
	if d.ValidFrom.IsZero() {
		d.ValidFrom = time.Now().UTC()
	}
	return s.repo.CreateDiscountRule(ctx, d)
}

func (s *catalogService) ListDiscountRules(ctx context.Context) ([]DiscountRule, error) {
	return s.repo.ListDiscountRules(ctx)
}
