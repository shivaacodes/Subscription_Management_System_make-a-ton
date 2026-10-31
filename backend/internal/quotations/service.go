package quotations

import (
	"context"
	"fmt"
	"time"
)

type SubscriptionCreator interface {
	CreateSubscriptionFromQuotation(ctx context.Context, q *Quotation, lines []QuotationLine) (string, error)
}

type Service interface {
	CreateQuotation(ctx context.Context, customerID, planID string, validityDays int, lines []QuotationLine) (*Quotation, error)
	GetQuotation(ctx context.Context, id string) (*Quotation, []QuotationLine, error)
	ListQuotations(ctx context.Context, customerID *string) ([]Quotation, error)
	ConvertQuotation(ctx context.Context, id string) (string, error)

	CreateTemplate(ctx context.Context, name string, description *string, validityDays int, planID *string, lines []QuotationTemplateLine) (*QuotationTemplate, error)
	ListTemplates(ctx context.Context) ([]QuotationTemplate, error)
}

type quotationService struct {
	repo       Repository
	subCreator SubscriptionCreator
}

func NewService(repo Repository, subCreator SubscriptionCreator) Service {
	return &quotationService{
		repo:       repo,
		subCreator: subCreator,
	}
}

func (s *quotationService) CreateQuotation(ctx context.Context, customerID, planID string, validityDays int, lines []QuotationLine) (*Quotation, error) {
	if validityDays <= 0 {
		validityDays = 30
	}

	var subtotal int64 = 0
	var totalTax int64 = 0

	for i := range lines {
		lineTotal := int64(lines[i].Quantity) * lines[i].UnitPriceCents
		taxAmount := (lineTotal * int64(lines[i].TaxRateBps)) / 10000
		lines[i].TotalCents = lineTotal + taxAmount

		subtotal += lineTotal
		totalTax += taxAmount
	}

	quoteNumber := fmt.Sprintf("QT-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

	q := &Quotation{
		QuotationNumber: quoteNumber,
		CustomerID:      customerID,
		PlanID:          planID,
		Status:          StatusSent,
		ValidUntil:      time.Now().UTC().AddDate(0, 0, validityDays),
		SubtotalCents:   subtotal,
		TaxCents:        totalTax,
		DiscountCents:   0,
		TotalCents:      subtotal + totalTax,
	}

	if err := s.repo.CreateQuotation(ctx, q, lines); err != nil {
		return nil, err
	}
	q.Lines = lines
	return q, nil
}

func (s *quotationService) GetQuotation(ctx context.Context, id string) (*Quotation, []QuotationLine, error) {
	return s.repo.GetQuotationByID(ctx, id)
}

func (s *quotationService) ListQuotations(ctx context.Context, customerID *string) ([]Quotation, error) {
	return s.repo.ListQuotations(ctx, customerID)
}

func (s *quotationService) ConvertQuotation(ctx context.Context, id string) (string, error) {
	q, lines, err := s.repo.GetQuotationByID(ctx, id)
	if err != nil {
		return "", err
	}

	if q.Status == StatusConverted {
		return "", fmt.Errorf("quotation is already converted")
	}

	if s.subCreator == nil {
		return "", fmt.Errorf("subscription creator not configured")
	}

	subID, err := s.subCreator.CreateSubscriptionFromQuotation(ctx, q, lines)
	if err != nil {
		return "", fmt.Errorf("failed converting quotation to subscription: %w", err)
	}

	_ = s.repo.UpdateStatus(ctx, id, StatusConverted)
	return subID, nil
}

func (s *quotationService) CreateTemplate(ctx context.Context, name string, description *string, validityDays int, planID *string, lines []QuotationTemplateLine) (*QuotationTemplate, error) {
	if validityDays <= 0 {
		validityDays = 30
	}
	t := &QuotationTemplate{
		Name:         name,
		Description:  description,
		ValidityDays: validityDays,
		PlanID:       planID,
		IsActive:     true,
	}
	if err := s.repo.CreateTemplate(ctx, t, lines); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *quotationService) ListTemplates(ctx context.Context) ([]QuotationTemplate, error) {
	return s.repo.ListTemplates(ctx)
}
