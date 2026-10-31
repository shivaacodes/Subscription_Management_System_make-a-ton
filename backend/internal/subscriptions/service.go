package subscriptions

import (
	"context"
	"fmt"
	"time"

	"github.com/subflow/backend/internal/plans"
	"github.com/subflow/backend/internal/quotations"
)

type PlanGetter interface {
	GetPlan(ctx context.Context, id string) (*plans.Plan, error)
	ComputeNextBillingDate(freq plans.Frequency, fromDate time.Time) time.Time
}

type Service interface {
	CreateSubscription(ctx context.Context, sub *Subscription, lines []SubscriptionLine, actorID *string) (*Subscription, error)
	GetSubscription(ctx context.Context, id string) (*Subscription, []SubscriptionLine, error)
	ListSubscriptions(ctx context.Context, customerID *string, status *string) ([]Subscription, error)

	TransitionStatus(ctx context.Context, id string, targetStatus SubscriptionStatus, actorID *string) error
	PauseSubscription(ctx context.Context, id string, actorID *string) error
	ResumeSubscription(ctx context.Context, id string, actorID *string) error
	CancelSubscription(ctx context.Context, id string, actorID *string) error

	ToggleAddon(ctx context.Context, subscriptionID, productID string, active bool, actorID *string) error
	GetTimeline(ctx context.Context, subscriptionID string) ([]SubscriptionEvent, error)

	CreateSubscriptionFromQuotation(ctx context.Context, q *quotations.Quotation, lines []quotations.QuotationLine) (string, error)
}

type subscriptionService struct {
	repo        Repository
	planService PlanGetter
}

func NewService(repo Repository, planService PlanGetter) Service {
	return &subscriptionService{
		repo:        repo,
		planService: planService,
	}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, sub *Subscription, lines []SubscriptionLine, actorID *string) (*Subscription, error) {
	plan, err := s.planService.GetPlan(ctx, sub.PlanID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan: %w", err)
	}

	if sub.StartDate.IsZero() {
		sub.StartDate = time.Now().UTC()
	}
	if sub.NextBillingDate.IsZero() {
		sub.NextBillingDate = s.planService.ComputeNextBillingDate(plan.Frequency, sub.StartDate)
	}
	if sub.SubscriptionNumber == "" {
		sub.SubscriptionNumber = fmt.Sprintf("SUB-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
	}
	if sub.Status == "" {
		sub.Status = StatusConfirmed
	}
	if sub.Currency == "" {
		sub.Currency = "INR"
	}
	if sub.PaymentTerms == "" {
		sub.PaymentTerms = "NET_30"
	}

	for i := range lines {
		lines[i].LineTotalCents = int64(lines[i].Quantity) * lines[i].UnitPriceCents
	}

	if err := s.repo.CreateSubscription(ctx, sub, lines); err != nil {
		return nil, err
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: sub.ID,
		EventType:      "SUBSCRIPTION_CREATED",
		ActorID:        actorID,
		Payload: map[string]any{
			"status":             sub.Status,
			"planId":             sub.PlanID,
			"subscriptionNumber": sub.SubscriptionNumber,
		},
	})

	return sub, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id string) (*Subscription, []SubscriptionLine, error) {
	return s.repo.GetSubscriptionByID(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, customerID *string, status *string) ([]Subscription, error) {
	return s.repo.ListSubscriptions(ctx, customerID, status)
}

func (s *subscriptionService) TransitionStatus(ctx context.Context, id string, targetStatus SubscriptionStatus, actorID *string) error {
	sub, _, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}

	oldStatus := sub.Status
	if oldStatus == targetStatus {
		return nil
	}

	// Validate Lifecycle transitions:
	// DRAFT -> QUOTATION -> CONFIRMED -> ACTIVE -> CLOSED
	// ACTIVE <-> PAUSED
	// Any non-closed -> CANCELLED
	switch targetStatus {
	case StatusActive:
		if oldStatus != StatusConfirmed && oldStatus != StatusPaused && oldStatus != StatusDraft {
			return fmt.Errorf("cannot activate subscription in state %s", oldStatus)
		}
	case StatusPaused:
		return s.PauseSubscription(ctx, id, actorID)
	case StatusCancelled:
		return s.CancelSubscription(ctx, id, actorID)
	case StatusClosed:
		if oldStatus == StatusClosed {
			return fmt.Errorf("subscription is already closed")
		}
	}

	if err := s.repo.UpdateStatus(ctx, id, targetStatus); err != nil {
		return err
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: id,
		EventType:      "STATUS_TRANSITION",
		ActorID:        actorID,
		Payload: map[string]any{
			"from": string(oldStatus),
			"to":   string(targetStatus),
		},
	})

	return nil
}

func (s *subscriptionService) PauseSubscription(ctx context.Context, id string, actorID *string) error {
	sub, _, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}

	if sub.Status != StatusActive {
		return fmt.Errorf("only ACTIVE subscriptions can be paused")
	}

	plan, err := s.planService.GetPlan(ctx, sub.PlanID)
	if err == nil && !plan.AllowPause {
		return fmt.Errorf("plan policy does not permit pausing this subscription")
	}

	if err := s.repo.UpdateStatus(ctx, id, StatusPaused); err != nil {
		return err
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: id,
		EventType:      "SUBSCRIPTION_PAUSED",
		ActorID:        actorID,
		Payload: map[string]any{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})

	return nil
}

func (s *subscriptionService) ResumeSubscription(ctx context.Context, id string, actorID *string) error {
	sub, _, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}

	if sub.Status != StatusPaused {
		return fmt.Errorf("only PAUSED subscriptions can be resumed")
	}

	if err := s.repo.UpdateStatus(ctx, id, StatusActive); err != nil {
		return err
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: id,
		EventType:      "SUBSCRIPTION_RESUMED",
		ActorID:        actorID,
		Payload: map[string]any{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})

	return nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, id string, actorID *string) error {
	sub, _, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}

	if sub.Status == StatusClosed || sub.Status == StatusCancelled {
		return fmt.Errorf("subscription is already terminated")
	}

	plan, err := s.planService.GetPlan(ctx, sub.PlanID)
	if err == nil && !plan.AllowEarlyClose {
		return fmt.Errorf("plan policy requires contract fulfillment until expiration")
	}

	if err := s.repo.UpdateStatus(ctx, id, StatusCancelled); err != nil {
		return err
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: id,
		EventType:      "SUBSCRIPTION_CANCELLED",
		ActorID:        actorID,
		Payload: map[string]any{
			"reason": "User or administrator cancellation request",
		},
	})

	return nil
}

func (s *subscriptionService) ToggleAddon(ctx context.Context, subscriptionID, productID string, active bool, actorID *string) error {
	if err := s.repo.UpdateAddonStatus(ctx, subscriptionID, productID, active); err != nil {
		return err
	}

	action := "ENABLED"
	if !active {
		action = "DISABLED"
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: subscriptionID,
		EventType:      "ADDON_MODIFIED",
		ActorID:        actorID,
		Payload: map[string]any{
			"productId": productID,
			"action":    action,
			"effective": "Next billing cycle",
		},
	})

	return nil
}

func (s *subscriptionService) GetTimeline(ctx context.Context, subscriptionID string) ([]SubscriptionEvent, error) {
	return s.repo.ListEvents(ctx, subscriptionID)
}

func (s *subscriptionService) CreateSubscriptionFromQuotation(ctx context.Context, q *quotations.Quotation, qLines []quotations.QuotationLine) (string, error) {
	plan, err := s.planService.GetPlan(ctx, q.PlanID)
	if err != nil {
		return "", fmt.Errorf("failed fetching plan: %w", err)
	}

	subNumber := fmt.Sprintf("SUB-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
	now := time.Now().UTC()
	nextBilling := s.planService.ComputeNextBillingDate(plan.Frequency, now)

	sub := &Subscription{
		SubscriptionNumber: subNumber,
		CustomerID:         q.CustomerID,
		PlanID:             q.PlanID,
		QuotationID:        &q.ID,
		Status:             StatusActive,
		StartDate:          now,
		NextBillingDate:    nextBilling,
		PaymentTerms:       "NET_30",
		Currency:           "INR",
	}

	var subLines []SubscriptionLine
	for _, ql := range qLines {
		subLines = append(subLines, SubscriptionLine{
			ProductID:      ql.ProductID,
			VariantID:      ql.VariantID,
			IsAddon:        false,
			Quantity:       ql.Quantity,
			UnitPriceCents: ql.UnitPriceCents,
			TaxRateBps:     ql.TaxRateBps,
			LineTotalCents: int64(ql.Quantity) * ql.UnitPriceCents,
		})
	}

	if err := s.repo.CreateSubscription(ctx, sub, subLines); err != nil {
		return "", err
	}

	_ = s.repo.RecordEvent(ctx, &SubscriptionEvent{
		SubscriptionID: sub.ID,
		EventType:      "CONVERTED_FROM_QUOTATION",
		Payload: map[string]any{
			"quotationNumber": q.QuotationNumber,
			"quotationId":     q.ID,
		},
	})

	return sub.ID, nil
}
