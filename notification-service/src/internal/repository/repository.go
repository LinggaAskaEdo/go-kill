package repository

import (
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/repository/notification"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	Notification notification.NotificationRepositoryItf
}

type Options struct {
	NotificationOpts notification.Options `yaml:"notification"`
}

func InitRepository(redis0 *redis.Client, mongo0 *mongo.Database) *Repository {
	return &Repository{
		Notification: notification.InitNotificationRepository(
			redis0,
			mongo0,
		),
	}
}
