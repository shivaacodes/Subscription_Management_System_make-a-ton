package plans

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidPlanPrice = errors.New("plan billing price must be non-negative")
	ErrInvalidQuantity  = errors.New("minimum quantity must be at least 1")
)

type Service interface {
	CreatePlan(ctx context.Context, p *Plan) error
	GetPlan(ctx context.Context, id string) (*Plan, error)
	ListPlans(ctx context.Context) ([]Plan, error)
	LinkProduct(ctx context.Context, pp *PlanProduct) error
	GetPlanProducts(ctx context.Context, planID string) ([]PlanProduct, error)
	ComputeNextBillingDate(freq Frequency, fromDate time.Time) time.Time
}

type planService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &planService{repo: repo}
}

func (s *planService) CreatePlan(ctx context.Context, p *Plan) error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return fmt.Errorf("plan name cannot be empty")
	}
	if p.BillingPriceCents < 0 {
		return ErrInvalidPlanPrice
	}
	if p.MinQuantity < 1 {
		p.MinQuantity = 1
	}
	if p.AvailableFrom.IsZero() {
		p.AvailableFrom = time.Now().UTC()
	}
	switch p.Frequency {
	case FrequencyDaily, FrequencyWeekly, FrequencyMonthly, FrequencyYearly:
		// valid
	default:
		return fmt.Errorf("invalid frequency: %s", p.Frequency)
	}
	return s.repo.CreatePlan(ctx, p)
}

func (s *planService) GetPlan(ctx context.Context, id string) (*Plan, error) {
	return s.repo.GetPlanByID(ctx, id)
}

func (s *planService) ListPlans(ctx context.Context) ([]Plan, error) {
	return s.repo.ListActivePlans(ctx)
}

func (s *planService) LinkProduct(ctx context.Context, pp *PlanProduct) error {
	return s.repo.AssociateProduct(ctx, pp)
}

func (s *planService) GetPlanProducts(ctx context.Context, planID string) ([]PlanProduct, error) {
	return s.repo.GetPlanProducts(ctx, planID)
}

// ComputeNextBillingDate correctly computes the subsequent billing date,
// maintaining month-end dates and handling leap years properly.
func (s *planService) ComputeNextBillingDate(freq Frequency, fromDate time.Time) time.Time {
	switch freq {
	case FrequencyDaily:
		return fromDate.AddDate(0, 0, 1)
	case FrequencyWeekly:
		return fromDate.AddDate(0, 0, 7)
	case FrequencyMonthly:
		return AddMonthsPreservingMonthEnd(fromDate, 1)
	case FrequencyYearly:
		return AddYearsPreservingLeap(fromDate, 1)
	default:
		return fromDate.AddDate(0, 1, 0)
	}
}

// AddMonthsPreservingMonthEnd adds n months and normalizes to month-end if needed
func AddMonthsPreservingMonthEnd(t time.Time, n int) time.Time {
	year, month, day := t.Date()
	isLastDay := isLastDayOfMonth(year, month, day)

	// Target year and month
	targetYear := year + (int(month)-1+n)/12
	targetMonth := time.Month(((int(month) - 1 + n) % 12) + 1)
	if targetMonth <= 0 {
		targetMonth += 12
		targetYear -= 1
	}

	maxDay := daysInMonth(targetYear, targetMonth)

	var targetDay int
	if isLastDay || day > maxDay {
		targetDay = maxDay
	} else {
		targetDay = day
	}

	return time.Date(targetYear, targetMonth, targetDay, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func AddYearsPreservingLeap(t time.Time, n int) time.Time {
	year, month, day := t.Date()
	targetYear := year + n

	if month == time.February && day == 29 && !isLeapYear(targetYear) {
		day = 28
	}

	return time.Date(targetYear, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func isLastDayOfMonth(year int, month time.Month, day int) bool {
	return day == daysInMonth(year, month)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
