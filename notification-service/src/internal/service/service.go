package service

import (
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/service/notification"
)

type Service struct {
	Notification notification.NotificationServiceItf
}

type Options struct {
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		Notification: notification.InitNotificationService(
			repository.Notification,
		),
	}
}
