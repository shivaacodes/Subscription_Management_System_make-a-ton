package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/subflow/backend/internal/billing"
)

type InvoiceManager interface {
	GetInvoice(ctx context.Context, id string) (*billing.Invoice, []billing.InvoiceLine, error)
	ConfirmInvoice(ctx context.Context, id string) error
}

type Service interface {
	RecordPayment(ctx context.Context, invoiceID string, amountCents int64, method string, notes *string) (*Payment, error)
	ListPayments(ctx context.Context, invoiceID *string, customerID *string) ([]Payment, error)
}

type paymentService struct {
	repo       Repository
	invManager InvoiceManager
}

func NewService(repo Repository, invManager InvoiceManager) Service {
	return &paymentService{
		repo:       repo,
		invManager: invManager,
	}
}

func (s *paymentService) RecordPayment(ctx context.Context, invoiceID string, amountCents int64, method string, notes *string) (*Payment, error) {
	if amountCents <= 0 {
		return nil, fmt.Errorf("payment amount must be greater than zero")
	}

	inv, _, err := s.invManager.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	alreadyPaid, err := s.repo.GetTotalPaidForInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed calculating existing payments: %w", err)
	}

	remaining := inv.TotalCents - alreadyPaid
	if amountCents > remaining {
		return nil, fmt.Errorf("payment amount (₹%.2f) exceeds outstanding balance (₹%.2f)",
			float64(amountCents)/100, float64(remaining)/100)
	}

	if method == "" {
		method = "BANK_TRANSFER"
	}

	paymentNumber := fmt.Sprintf("PAY-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

	p := &Payment{
		InvoiceID:     invoiceID,
		PaymentNumber: paymentNumber,
		AmountCents:   amountCents,
		PaymentMethod: method,
		PaymentDate:   time.Now().UTC(),
		Notes:         notes,
	}

	if err := s.repo.RecordPayment(ctx, p); err != nil {
		return nil, fmt.Errorf("failed recording payment: %w", err)
	}

	return p, nil
}

func (s *paymentService) ListPayments(ctx context.Context, invoiceID *string, customerID *string) ([]Payment, error) {
	if invoiceID != nil && *invoiceID != "" {
		return s.repo.ListPaymentsByInvoiceID(ctx, *invoiceID)
	}
	return s.repo.ListAllPayments(ctx, customerID)
}
