package reporting

import (
	"context"
	"time"
)

type DashboardMetrics struct {
	ActiveSubscriptionsCount int64 `json:"activeSubscriptionsCount"`
	MRRCents                 int64 `json:"mrrCents"`
	InvoicedRevenueCents     int64 `json:"invoicedRevenueCents"`
	CollectedPaymentsCents   int64 `json:"collectedPaymentsCents"`
	OutstandingReceivables   int64 `json:"outstandingReceivablesCents"`
	OverdueInvoicesCount     int64 `json:"overdueInvoicesCount"`
}

type Repository interface {
	GetDashboardMetrics(ctx context.Context, asOf time.Time) (*DashboardMetrics, error)
}
