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
	CustomerEmail   string          `json:"customerEmail,omitempty"`
	PlanID          string          `json:"planId"`
	PlanName        string          `json:"planName,omitempty"`
	Status          QuotationStatus `json:"status"`
	ValidUntil      time.Time       `json:"validUntil"`
	SubtotalCents   int64           `json:"subtotalCents"`
	TaxCents        int64           `json:"taxCents"`
	DiscountCents   int64           `json:"discountCents"`
	TotalCents      int64           `json:"totalCents"`
	Lines           []QuotationLine `json:"lines,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type QuotationLine struct {
	ID             string  `json:"id"`
	QuotationID    string  `json:"quotationId"`
	ProductID      string  `json:"productId"`
	ProductName    string  `json:"productName,omitempty"`
	VariantID      *string `json:"variantId,omitempty"`
	VariantName    *string `json:"variantName,omitempty"`
	Quantity       int     `json:"quantity"`
	UnitPriceCents int64   `json:"unitPriceCents"`
	TaxRateBps     int     `json:"taxRateBps"`
	TotalCents     int64   `json:"totalCents"`
}

type QuotationTemplate struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Description  *string                 `json:"description,omitempty"`
	ValidityDays int                     `json:"validityDays"`
	PlanID       *string                 `json:"planId,omitempty"`
	IsActive     bool                    `json:"isActive"`
	Lines        []QuotationTemplateLine `json:"lines,omitempty"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

type QuotationTemplateLine struct {
	ID             string  `json:"id"`
	TemplateID     string  `json:"templateId"`
	ProductID      string  `json:"productId"`
	VariantID      *string `json:"variantId,omitempty"`
	Quantity       int     `json:"quantity"`
	UnitPriceCents int64   `json:"unitPriceCents"`
}

type Repository interface {
	CreateQuotation(ctx context.Context, q *Quotation, lines []QuotationLine) error
	GetQuotationByID(ctx context.Context, id string) (*Quotation, []QuotationLine, error)
	ListQuotations(ctx context.Context, customerID *string) ([]Quotation, error)
	UpdateStatus(ctx context.Context, id string, status QuotationStatus) error

	CreateTemplate(ctx context.Context, t *QuotationTemplate, lines []QuotationTemplateLine) error
	ListTemplates(ctx context.Context) ([]QuotationTemplate, error)
}
