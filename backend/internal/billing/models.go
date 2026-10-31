package billing

import (
	"context"
	"time"
)

type InvoiceStatus string

const (
	InvoiceDraft     InvoiceStatus = "DRAFT"
	InvoiceConfirmed InvoiceStatus = "CONFIRMED"
	InvoicePaid      InvoiceStatus = "PAID"
	InvoiceCancelled InvoiceStatus = "CANCELLED"
)

type Invoice struct {
	ID                 string        `json:"id"`
	InvoiceNumber      string        `json:"invoiceNumber"`
	SubscriptionID     string        `json:"subscriptionId"`
	CustomerID         string        `json:"customerId"`
	BillingRunID       *string       `json:"billingRunId,omitempty"`
	Status             InvoiceStatus `json:"status"`
	BillingPeriodStart time.Time     `json:"billingPeriodStart"`
	BillingPeriodEnd   time.Time     `json:"billingPeriodEnd"`
	Currency           string        `json:"currency"`
	SubtotalCents      int64         `json:"subtotalCents"`
	TaxCents           int64         `json:"taxCents"`
	DiscountCents      int64         `json:"discountCents"`
	TotalCents         int64         `json:"totalCents"`
	DueDate            time.Time     `json:"dueDate"`
	PaidAt             *time.Time    `json:"paidAt,omitempty"`
	CreatedAt          time.Time     `json:"createdAt"`
	UpdatedAt          time.Time     `json:"updatedAt"`
}

type InvoiceLine struct {
	ID             string    `json:"id"`
	InvoiceID      string    `json:"invoiceId"`
	ProductID      string    `json:"productId"`
	VariantID      *string   `json:"variantId,omitempty"`
	Description    string    `json:"description"`
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unitPriceCents"`
	TaxCents       int64     `json:"taxCents"`
	DiscountCents  int64     `json:"discountCents"`
	TotalCents     int64     `json:"totalCents"`
	CreatedAt      time.Time `json:"createdAt"`
}

// ExplainableBreakdown represents the authoritative billing calculation contract
type ExplainableBreakdown struct {
	BasePriceCents      int64 `json:"basePriceCents"`
	VariantAdjustment   int64 `json:"variantAdjustmentCents"`
	EffectiveUnitPrice  int64 `json:"effectiveUnitPriceCents"`
	Quantity            int   `json:"quantity"`
	GrossAmountCents    int64 `json:"grossAmountCents"`
	DiscountAmountCents int64 `json:"discountAmountCents"`
	TaxAmountCents      int64 `json:"taxAmountCents"`
	FinalTotalCents     int64 `json:"finalTotalCents"`
}

type Repository interface {
	CreateInvoice(ctx context.Context, inv *Invoice, lines []InvoiceLine) error
	GetInvoiceByID(ctx context.Context, id string) (*Invoice, []InvoiceLine, error)
	ListCustomerInvoices(ctx context.Context, customerID string) ([]Invoice, error)
	InvoiceExistsForPeriod(ctx context.Context, subscriptionID string, start, end time.Time) (bool, error)
}
