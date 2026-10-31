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
	ID                      string        `json:"id"`
	InvoiceNumber           string        `json:"invoiceNumber"`
	SubscriptionID          string        `json:"subscriptionId"`
	SubscriptionNumber      string        `json:"subscriptionNumber,omitempty"`
	CustomerID              string        `json:"customerId"`
	CustomerEmail           string        `json:"customerEmail,omitempty"`
	BillingRunID            *string       `json:"billingRunId,omitempty"`
	Status                  InvoiceStatus `json:"status"`
	BillingPeriodStart      time.Time     `json:"billingPeriodStart"`
	BillingPeriodEnd        time.Time     `json:"billingPeriodEnd"`
	Currency                string        `json:"currency"`
	SubtotalCents           int64         `json:"subtotalCents"`
	TaxCents                int64         `json:"taxCents"`
	DiscountCents           int64         `json:"discountCents"`
	TotalCents              int64         `json:"totalCents"`
	PaidAmountCents         int64         `json:"paidAmountCents"`
	OutstandingBalanceCents int64         `json:"outstandingBalanceCents"`
	DueDate                 time.Time     `json:"dueDate"`
	PaidAt                  *time.Time    `json:"paidAt,omitempty"`
	Lines                   []InvoiceLine `json:"lines,omitempty"`
	CreatedAt               time.Time     `json:"createdAt"`
	UpdatedAt               time.Time     `json:"updatedAt"`
}

type InvoiceLine struct {
	ID             string    `json:"id"`
	InvoiceID      string    `json:"invoiceId"`
	ProductID      string    `json:"productId"`
	ProductName    string    `json:"productName,omitempty"`
	VariantID      *string   `json:"variantId,omitempty"`
	VariantName    *string   `json:"variantName,omitempty"`
	Description    string    `json:"description"`
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unitPriceCents"`
	TaxCents       int64     `json:"taxCents"`
	DiscountCents  int64     `json:"discountCents"`
	TotalCents     int64     `json:"totalCents"`
	CreatedAt      time.Time `json:"createdAt"`
}

type BillingRun struct {
	ID                     string           `json:"id"`
	RunDate                time.Time        `json:"runDate"`
	Status                 string           `json:"status"`
	SubscriptionsEvaluated int              `json:"subscriptionsEvaluated"`
	InvoicesGenerated      int              `json:"invoicesGenerated"`
	FailuresCount          int              `json:"failuresCount"`
	ErrorSummary           *string          `json:"errorSummary,omitempty"`
	Attempts               []BillingAttempt `json:"attempts,omitempty"`
	CreatedAt              time.Time        `json:"createdAt"`
	CompletedAt            *time.Time       `json:"completedAt,omitempty"`
}

type BillingAttempt struct {
	ID                 string    `json:"id"`
	BillingRunID       string    `json:"billingRunId"`
	SubscriptionID     string    `json:"subscriptionId"`
	SubscriptionNumber string    `json:"subscriptionNumber,omitempty"`
	Status             string    `json:"status"` // SUCCESS, FAILED, RETRYABLE
	ErrorMessage       *string   `json:"errorMessage,omitempty"`
	AttemptNumber      int       `json:"attemptNumber"`
	CreatedAt          time.Time `json:"createdAt"`
}

type ReliabilityMetrics struct {
	SubscriptionsDueCount int              `json:"subscriptionsDueCount"`
	RecentRuns            []BillingRun     `json:"recentRuns"`
	RecentAttempts        []BillingAttempt `json:"recentAttempts"`
	DuplicatePrevented    int              `json:"duplicatePreventedCount"`
}

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
	ListInvoices(ctx context.Context, customerID *string, status *string) ([]Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, id string, status InvoiceStatus) error
	InvoiceExistsForPeriod(ctx context.Context, subscriptionID string, start, end time.Time) (bool, error)

	CreateBillingRun(ctx context.Context, run *BillingRun) error
	CompleteBillingRun(ctx context.Context, id string, evaluated, generated, failures int, errorSummary *string) error
	RecordBillingAttempt(ctx context.Context, attempt *BillingAttempt) error
	ListRecentBillingRuns(ctx context.Context, limit int) ([]BillingRun, error)
	ListRecentAttempts(ctx context.Context, limit int) ([]BillingAttempt, error)
}
