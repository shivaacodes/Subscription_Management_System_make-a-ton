package reporting

import (
	"context"
	"time"
)

type MonthlyTrendPoint struct {
	Month                string `json:"month"`
	InvoicedRevenueCents int64  `json:"invoicedRevenueCents"`
	CollectedCents       int64  `json:"collectedCents"`
}

type DashboardMetrics struct {
	ActiveSubscriptionsCount int64               `json:"activeSubscriptionsCount"`
	MRRCents                 int64               `json:"mrrCents"`
	InvoicedRevenueCents     int64               `json:"invoicedRevenueCents"`
	CollectedPaymentsCents   int64               `json:"collectedPaymentsCents"`
	OutstandingReceivables   int64               `json:"outstandingReceivablesCents"`
	OverdueInvoicesCount     int64               `json:"overdueInvoicesCount"`
	Trend                    []MonthlyTrendPoint `json:"trend"`
}

type Repository interface {
	GetDashboardMetrics(ctx context.Context, asOf time.Time) (*DashboardMetrics, error)
}
