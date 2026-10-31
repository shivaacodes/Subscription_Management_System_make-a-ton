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
	CustomerEmail      string             `json:"customerEmail,omitempty"`
	PlanID             string             `json:"planId"`
	PlanName           string             `json:"planName,omitempty"`
	QuotationID        *string            `json:"quotationId,omitempty"`
	Status             SubscriptionStatus `json:"status"`
	StartDate          time.Time          `json:"startDate"`
	ExpiryDate         *time.Time         `json:"expiryDate,omitempty"`
	NextBillingDate    time.Time          `json:"nextBillingDate"`
	PaymentTerms       string             `json:"paymentTerms"`
	Currency           string             `json:"currency"`
	Lines              []SubscriptionLine `json:"lines,omitempty"`
	TotalAmountCents   int64              `json:"totalAmountCents,omitempty"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

type SubscriptionLine struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscriptionId"`
	ProductID      string    `json:"productId"`
	ProductName    string    `json:"productName,omitempty"`
	VariantID      *string   `json:"variantId,omitempty"`
	VariantName    *string   `json:"variantName,omitempty"`
	IsAddon        bool      `json:"isAddon"`
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unitPriceCents"`
	TaxRateBps     int       `json:"taxRateBps"`
	LineTotalCents int64     `json:"lineTotalCents"`
	CreatedAt      time.Time `json:"createdAt"`
}

type SubscriptionEvent struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscriptionId"`
	EventType      string    `json:"eventType"`
	ActorID        *string   `json:"actorId,omitempty"`
	ActorEmail     string    `json:"actorEmail,omitempty"`
	Payload        any       `json:"payload"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Repository interface {
	CreateSubscription(ctx context.Context, sub *Subscription, lines []SubscriptionLine) error
	GetSubscriptionByID(ctx context.Context, id string) (*Subscription, []SubscriptionLine, error)
	ListSubscriptions(ctx context.Context, customerID *string, status *string) ([]Subscription, error)
	GetSubscriptionsDueForBilling(ctx context.Context, cutoffDate time.Time, limit int) ([]Subscription, error)
	UpdateStatus(ctx context.Context, id string, status SubscriptionStatus) error
	UpdateNextBillingDate(ctx context.Context, id string, nextBillingDate time.Time) error
	UpdateAddonStatus(ctx context.Context, subscriptionID, productID string, active bool) error

	RecordEvent(ctx context.Context, event *SubscriptionEvent) error
	ListEvents(ctx context.Context, subscriptionID string) ([]SubscriptionEvent, error)
}
