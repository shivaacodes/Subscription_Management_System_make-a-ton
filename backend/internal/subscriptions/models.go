package subscriptions

import (
	"context"
	"time"
)

type SubscriptionStatus string

const (
	StatusDraft     SubscriptionStatus = "DRAFT"
	StatusQuotation SubscriptionStatus = "QUOTATION"
	StatusConfirmed SubscriptionStatus = "CONFIRMED"
	StatusActive    SubscriptionStatus = "ACTIVE"
	StatusPaused    SubscriptionStatus = "PAUSED"
	StatusCancelled SubscriptionStatus = "CANCELLED"
	StatusClosed    SubscriptionStatus = "CLOSED"
)

type Subscription struct {
	ID                 string             `json:"id"`
	SubscriptionNumber string             `json:"subscriptionNumber"`
	CustomerID         string             `json:"customerId"`
	PlanID             string             `json:"planId"`
	QuotationID        *string            `json:"quotationId,omitempty"`
	Status             SubscriptionStatus `json:"status"`
	StartDate          time.Time          `json:"startDate"`
	ExpiryDate         *time.Time         `json:"expiryDate,omitempty"`
	NextBillingDate    time.Time          `json:"nextBillingDate"`
	PaymentTerms       string             `json:"paymentTerms"`
	Currency           string             `json:"currency"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

type SubscriptionLine struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscriptionId"`
	ProductID      string    `json:"productId"`
	VariantID      *string   `json:"variantId,omitempty"`
	IsAddon        bool      `json:"isAddon"`
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unitPriceCents"`
	TaxRateBps     int       `json:"taxRateBps"`
	LineTotalCents int64     `json:"lineTotalCents"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Repository interface {
	CreateSubscription(ctx context.Context, sub *Subscription, lines []SubscriptionLine) error
	GetSubscriptionByID(ctx context.Context, id string) (*Subscription, []SubscriptionLine, error)
	GetSubscriptionsDueForBilling(ctx context.Context, cutoffDate time.Time, limit int) ([]Subscription, error)
	UpdateStatus(ctx context.Context, id string, status SubscriptionStatus) error
	UpdateNextBillingDate(ctx context.Context, id string, nextBillingDate time.Time) error
}
