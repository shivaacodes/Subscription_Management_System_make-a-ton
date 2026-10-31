package billing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/subflow/backend/internal/plans"
	"github.com/subflow/backend/internal/subscriptions"
)

type SubscriptionManager interface {
	GetSubscriptionsDueForBilling(ctx context.Context, cutoffDate time.Time, limit int) ([]subscriptions.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id string) (*subscriptions.Subscription, []subscriptions.SubscriptionLine, error)
	UpdateNextBillingDate(ctx context.Context, id string, nextBillingDate time.Time) error
}

type PlanProvider interface {
	GetPlan(ctx context.Context, id string) (*plans.Plan, error)
	ComputeNextBillingDate(freq plans.Frequency, fromDate time.Time) time.Time
}

type Service interface {
	RunAutomatedBilling(ctx context.Context) (*BillingRun, error)
	GetReliabilityMetrics(ctx context.Context) (*ReliabilityMetrics, error)
	GetInvoice(ctx context.Context, id string) (*Invoice, []InvoiceLine, error)
	ListInvoices(ctx context.Context, customerID *string, status *string) ([]Invoice, error)
	ConfirmInvoice(ctx context.Context, id string) error
	PreviewCalculation(ctx context.Context, basePrice, variantAdj int64, qty int, discount int64, taxBps int) ExplainableBreakdown
}

type billingService struct {
	repo        Repository
	subManager  SubscriptionManager
	planService PlanProvider
}

func NewService(repo Repository, subManager SubscriptionManager, planService PlanProvider) Service {
	return &billingService{
		repo:        repo,
		subManager:  subManager,
		planService: planService,
	}
}

func (s *billingService) PreviewCalculation(_ context.Context, basePrice, variantAdj int64, qty int, discount int64, taxBps int) ExplainableBreakdown {
	return CalculateLineBreakdown(basePrice, variantAdj, qty, discount, taxBps)
}

func (s *billingService) RunAutomatedBilling(ctx context.Context) (*BillingRun, error) {
	now := time.Now().UTC()
	run := &BillingRun{
		RunDate:                now,
		Status:                 "IN_PROGRESS",
		SubscriptionsEvaluated: 0,
		InvoicesGenerated:      0,
		FailuresCount:          0,
	}

	if err := s.repo.CreateBillingRun(ctx, run); err != nil {
		return nil, fmt.Errorf("failed to create billing run record: %w", err)
	}

	dueSubs, err := s.subManager.GetSubscriptionsDueForBilling(ctx, now, 500)
	if err != nil {
		summary := fmt.Sprintf("failed fetching due subscriptions: %v", err)
		_ = s.repo.CompleteBillingRun(ctx, run.ID, 0, 0, 1, &summary)
		return run, err
	}

	run.SubscriptionsEvaluated = len(dueSubs)

	for _, sub := range dueSubs {
		periodStart := sub.StartDate
		periodEnd := sub.NextBillingDate

		// Idempotency: verify if invoice already exists for this exact period
		exists, err := s.repo.InvoiceExistsForPeriod(ctx, sub.ID, periodStart, periodEnd)
		if err != nil {
			slog.Error("error checking existing invoice for period", "subId", sub.ID, "err", err)
			continue
		}

		if exists {
			// Duplicate prevented safely!
			_ = s.repo.RecordBillingAttempt(ctx, &BillingAttempt{
				BillingRunID:   run.ID,
				SubscriptionID: sub.ID,
				Status:         "SUCCESS",
				ErrorMessage:   ptrString("Skipped: invoice already generated for this billing cycle"),
				AttemptNumber:  1,
			})
			continue
		}

		plan, err := s.planService.GetPlan(ctx, sub.PlanID)
		if err != nil {
			errMsg := fmt.Sprintf("plan error: %v", err)
			_ = s.repo.RecordBillingAttempt(ctx, &BillingAttempt{
				BillingRunID:   run.ID,
				SubscriptionID: sub.ID,
				Status:         "FAILED",
				ErrorMessage:   &errMsg,
				AttemptNumber:  1,
			})
			run.FailuresCount++
			continue
		}

		// Calculate lines
		var invLines []InvoiceLine
		var subtotal int64 = 0
		var taxTotal int64 = 0

		for _, sl := range sub.Lines {
			breakdown := CalculateLineBreakdown(sl.UnitPriceCents, 0, sl.Quantity, 0, sl.TaxRateBps)
			desc := fmt.Sprintf("%s x %d", sl.ProductName, sl.Quantity)
			if desc == " x 1" || sl.ProductName == "" {
				desc = fmt.Sprintf("Subscription billing item (qty: %d)", sl.Quantity)
			}

			invLines = append(invLines, InvoiceLine{
				ProductID:      sl.ProductID,
				VariantID:      sl.VariantID,
				Description:    desc,
				Quantity:       sl.Quantity,
				UnitPriceCents: sl.UnitPriceCents,
				TaxCents:       breakdown.TaxAmountCents,
				DiscountCents:  0,
				TotalCents:     breakdown.FinalTotalCents,
			})

			subtotal += breakdown.GrossAmountCents
			taxTotal += breakdown.TaxAmountCents
		}

		invNumber := fmt.Sprintf("INV-%d-%04d", now.Year(), now.UnixNano()%100000)
		inv := &Invoice{
			InvoiceNumber:      invNumber,
			SubscriptionID:     sub.ID,
			CustomerID:         sub.CustomerID,
			BillingRunID:       &run.ID,
			Status:             InvoiceConfirmed,
			BillingPeriodStart: periodStart,
			BillingPeriodEnd:   periodEnd,
			Currency:           sub.Currency,
			SubtotalCents:      subtotal,
			TaxCents:           taxTotal,
			DiscountCents:      0,
			TotalCents:         subtotal + taxTotal,
			DueDate:            now.AddDate(0, 0, 30),
		}

		if err := s.repo.CreateInvoice(ctx, inv, invLines); err != nil {
			errMsg := fmt.Sprintf("failed creating invoice: %v", err)
			_ = s.repo.RecordBillingAttempt(ctx, &BillingAttempt{
				BillingRunID:   run.ID,
				SubscriptionID: sub.ID,
				Status:         "FAILED",
				ErrorMessage:   &errMsg,
				AttemptNumber:  1,
			})
			run.FailuresCount++
			continue
		}

		// Advance subscription next_billing_date
		nextDate := s.planService.ComputeNextBillingDate(plan.Frequency, sub.NextBillingDate)
		_ = s.subManager.UpdateNextBillingDate(ctx, sub.ID, nextDate)

		_ = s.repo.RecordBillingAttempt(ctx, &BillingAttempt{
			BillingRunID:   run.ID,
			SubscriptionID: sub.ID,
			Status:         "SUCCESS",
			AttemptNumber:  1,
		})
		run.InvoicesGenerated++
	}

	_ = s.repo.CompleteBillingRun(ctx, run.ID, run.SubscriptionsEvaluated, run.InvoicesGenerated, run.FailuresCount, nil)
	return run, nil
}

func (s *billingService) GetReliabilityMetrics(ctx context.Context) (*ReliabilityMetrics, error) {
	now := time.Now().UTC()
	dueSubs, _ := s.subManager.GetSubscriptionsDueForBilling(ctx, now, 1000)
	runs, _ := s.repo.ListRecentBillingRuns(ctx, 10)
	attempts, _ := s.repo.ListRecentAttempts(ctx, 20)

	var dupPrevented = 0
	for _, a := range attempts {
		if a.ErrorMessage != nil && *a.ErrorMessage == "Skipped: invoice already generated for this billing cycle" {
			dupPrevented++
		}
	}

	return &ReliabilityMetrics{
		SubscriptionsDueCount: len(dueSubs),
		RecentRuns:            runs,
		RecentAttempts:        attempts,
		DuplicatePrevented:    dupPrevented,
	}, nil
}

func (s *billingService) GetInvoice(ctx context.Context, id string) (*Invoice, []InvoiceLine, error) {
	return s.repo.GetInvoiceByID(ctx, id)
}

func (s *billingService) ListInvoices(ctx context.Context, customerID *string, status *string) ([]Invoice, error) {
	return s.repo.ListInvoices(ctx, customerID, status)
}

func (s *billingService) ConfirmInvoice(ctx context.Context, id string) error {
	return s.repo.UpdateInvoiceStatus(ctx, id, InvoiceConfirmed)
}

func ptrString(s string) *string {
	return &s
}
