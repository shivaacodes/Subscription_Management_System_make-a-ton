package plans

import (
	"context"
	"time"
)

type Frequency string

const (
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyYearly  Frequency = "YEARLY"
)

type Plan struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Frequency         Frequency  `json:"frequency"`
	BillingPriceCents int64      `json:"billingPriceCents"`
	MinQuantity       int        `json:"minQuantity"`
	AvailableFrom     time.Time  `json:"availableFrom"`
	AvailableTo       *time.Time `json:"availableTo,omitempty"`
	AutoClose         bool       `json:"autoClose"`
	AllowEarlyClose   bool       `json:"allowEarlyClose"`
	AllowPause        bool       `json:"allowPause"`
	AllowRenewal      bool       `json:"allowRenewal"`
	IsActive          bool       `json:"isActive"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

type PlanProduct struct {
	ID         string `json:"id"`
	PlanID     string `json:"planId"`
	ProductID  string `json:"productId"`
	IsAddon    bool   `json:"isAddon"`
	IsRequired bool   `json:"isRequired"`
}

type Repository interface {
	CreatePlan(ctx context.Context, p *Plan) error
	GetPlanByID(ctx context.Context, id string) (*Plan, error)
	ListActivePlans(ctx context.Context) ([]Plan, error)
	AssociateProduct(ctx context.Context, pp *PlanProduct) error
	GetPlanProducts(ctx context.Context, planID string) ([]PlanProduct, error)
}
