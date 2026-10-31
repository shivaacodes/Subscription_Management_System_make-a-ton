package reporting

import (
	"context"
	"time"
)

type Service interface {
	GetDashboard(ctx context.Context) (*DashboardMetrics, error)
}

type reportingService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &reportingService{repo: repo}
}

func (s *reportingService) GetDashboard(ctx context.Context) (*DashboardMetrics, error) {
	return s.repo.GetDashboardMetrics(ctx, time.Now().UTC())
}
