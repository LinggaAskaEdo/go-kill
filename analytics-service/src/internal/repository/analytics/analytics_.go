package analytics

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"
)

func (r *analyticsRepository) UpdateOrderAnalytics(ctx context.Context, event dto.OrderEvent) error {
	date := time.Date(event.Timestamp.Year(), event.Timestamp.Month(), event.Timestamp.Day(), 0, 0, 0, 0, time.UTC)

	// Update order mongo
	err := r.updateOrderMongo(ctx, date, event)
	if err != nil {
		return err
	}

	// Get data from mongo
	resultData, err := r.getOrderMongo(ctx, date)
	if err != nil {
		return err
	}

	// Update redis
	err = r.updateOrderCache(ctx, date, resultData)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) UpdateProductAnalytics(ctx context.Context, event dto.OrderEvent) error {
	date := time.Date(event.Timestamp.Year(), event.Timestamp.Month(), event.Timestamp.Day(), 0, 0, 0, 0, time.UTC)

	// Update product mongo
	err := r.updateProductMongo(ctx, date, event)
	if err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepository) UpdateCancellationMetrics(ctx context.Context, event dto.OrderEvent) error {
	date := time.Date(event.Timestamp.Year(), event.Timestamp.Month(), event.Timestamp.Day(), 0, 0, 0, 0, time.UTC)

	// Update product mongo
	err := r.updateCancellationMongo(ctx, date)
	if err != nil {
		return err
	}

	// Get data from mongo
	resultData, err := r.getOrderMongo(ctx, date)
	if err != nil {
		return err
	}

	// Update redis
	err = r.updateOrderCache(ctx, date, resultData)
	if err != nil {
		return err
	}

	return nil
}
