package analytics

import (
	"context"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/repository/analytics"
)

type AnalyticsServiceItf interface {
	UpdateOrderAnalytics(ctx context.Context, event dto.OrderEvent) error
	UpdateProductAnalytics(ctx context.Context, event dto.OrderEvent) error
	UpdateCancellationMetrics(ctx context.Context, event dto.OrderEvent) error
}

type analyticsService struct {
	analyticsRepository analytics.AnalyticsRepositoryItf
}

type Options struct {
}

func InitAnalyticsService(analyticsRepository analytics.AnalyticsRepositoryItf) AnalyticsServiceItf {
	return &analyticsService{
		analyticsRepository: analyticsRepository,
	}
}
