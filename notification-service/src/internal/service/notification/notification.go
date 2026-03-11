package notification

import (
	"context"

	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/repository/notification"
)

type NotificationServiceItf interface {
	GetUserPreference(ctx context.Context, userID string) (dto.NotificationPreferences, error)
	CheckRateLimit(ctx context.Context, userID string) bool
	SendOrderConfirmation(ctx context.Context, event dto.OrderEvent) error
	SendOrderUpdate(ctx context.Context, event dto.OrderEvent) error
	SendOrderCancellation(ctx context.Context, event dto.OrderEvent) error
}

type notificationService struct {
	notificationRepository notification.NotificationRepositoryItf
}

type Options struct {
}

func InitNotificationService(notificationRepository notification.NotificationRepositoryItf) NotificationServiceItf {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}
