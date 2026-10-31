package quotations

import (
	"context"
	"time"
)

type QuotationStatus string

const (
	StatusDraft     QuotationStatus = "DRAFT"
	StatusSent      QuotationStatus = "SENT"
	StatusAccepted  QuotationStatus = "ACCEPTED"
	StatusExpired   QuotationStatus = "EXPIRED"
	StatusConverted QuotationStatus = "CONVERTED"
)

type Quotation struct {
	ID              string          `json:"id"`
	QuotationNumber string          `json:"quotationNumber"`
	CustomerID      string          `json:"customerId"`
	PlanID          string          `json:"planId"`
	Status          QuotationStatus `json:"status"`
	ValidUntil      time.Time       `json:"validUntil"`
	SubtotalCents   int64           `json:"subtotalCents"`
	TaxCents        int64           `json:"taxCents"`
	DiscountCents   int64           `json:"discountCents"`
	TotalCents      int64           `json:"totalCents"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type QuotationLine struct {
	ID             string  `json:"id"`
	QuotationID    string  `json:"quotationId"`
	ProductID      string  `json:"productId"`
	VariantID      *string `json:"variantId,omitempty"`
	Quantity       int     `json:"quantity"`
	UnitPriceCents int64   `json:"unitPriceCents"`
	TaxRateBps     int     `json:"taxRateBps"`
	TotalCents     int64   `json:"totalCents"`
}

type Repository interface {
	CreateQuotation(ctx context.Context, q *Quotation, lines []QuotationLine) error
	GetQuotationByID(ctx context.Context, id string) (*Quotation, []QuotationLine, error)
	ListCustomerQuotations(ctx context.Context, customerID string) ([]Quotation, error)
}
