package payments

import (
	"context"
	"time"
)

type Payment struct {
	ID            string    `json:"id"`
	InvoiceID     string    `json:"invoiceId"`
	PaymentNumber string    `json:"paymentNumber"`
	AmountCents   int64     `json:"amountCents"`
	PaymentMethod string    `json:"paymentMethod"`
	PaymentDate   time.Time `json:"paymentDate"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Repository interface {
	RecordPayment(ctx context.Context, p *Payment) error
	ListPaymentsByInvoiceID(ctx context.Context, invoiceID string) ([]Payment, error)
	GetTotalPaidForInvoice(ctx context.Context, invoiceID string) (int64, error)
}
