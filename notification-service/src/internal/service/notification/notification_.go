package notification

import (
	"context"

	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/model/dto"
)

func (s *notificationService) GetUserPreference(ctx context.Context, userID string) (dto.NotificationPreferences, error) {
	return s.notificationRepository.GetUserPreference(ctx, userID)
}

func (s *notificationService) CheckRateLimit(ctx context.Context, userID string) bool {
	return s.notificationRepository.CheckRateLimit(ctx, userID)
}

func (s *notificationService) SendOrderConfirmation(ctx context.Context, event dto.OrderEvent) error {
	return s.notificationRepository.SendOrderConfirmation(ctx, event)
}

func (s *notificationService) SendOrderUpdate(ctx context.Context, event dto.OrderEvent) error {
	return s.notificationRepository.SendOrderConfirmation(ctx, event)
}

func (s *notificationService) SendOrderCancellation(ctx context.Context, event dto.OrderEvent) error {
	return s.notificationRepository.SendOrderCancellation(ctx, event)
}
