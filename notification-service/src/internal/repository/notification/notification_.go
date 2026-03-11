package notification

import (
	"context"

	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/model/dto"
)

func (r *notificationRepository) GetUserPreference(ctx context.Context, userID string) (dto.NotificationPreferences, error) {
	return r.getUserPreferenceMongo(ctx, userID)
}

func (r *notificationRepository) CheckRateLimit(ctx context.Context, userID string) bool {
	return r.checkRateLimitCache(ctx, userID)
}

func (r *notificationRepository) SendOrderConfirmation(ctx context.Context, event dto.OrderEvent) error {
	return r.sendOrderConfirmationMongo(ctx, event)
}

func (r *notificationRepository) SendOrderUpdate(ctx context.Context, event dto.OrderEvent) error {
	return r.sendOrderUpdateMongo(ctx, event)
}

func (r *notificationRepository) SendOrderCancellation(ctx context.Context, event dto.OrderEvent) error {
	return r.sendOrderCancellationMongo(ctx, event)
}
