package analytics

import (
	"context"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"
)

func (s *analyticsService) UpdateOrderAnalytics(ctx context.Context, event dto.OrderEvent) error {
	return s.analyticsRepository.UpdateOrderAnalytics(ctx, event)
}

func (s *analyticsService) UpdateProductAnalytics(ctx context.Context, event dto.OrderEvent) error {
	return s.analyticsRepository.UpdateProductAnalytics(ctx, event)
}

func (s *analyticsService) UpdateCancellationMetrics(ctx context.Context, event dto.OrderEvent) error {
	return s.analyticsRepository.UpdateCancellationMetrics(ctx, event)
}
