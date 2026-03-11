package notification

import (
	"context"
	"fmt"

	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/model/dto"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationRepositoryItf interface {
	GetUserPreference(ctx context.Context, userID string) (dto.NotificationPreferences, error)
	CheckRateLimit(ctx context.Context, userID string) bool
	SendOrderConfirmation(ctx context.Context, event dto.OrderEvent) error
	SendOrderUpdate(ctx context.Context, event dto.OrderEvent) error
	SendOrderCancellation(ctx context.Context, event dto.OrderEvent) error
}

type notificationRepository struct {
	redis0 *redis.Client
	mongo0 *mongo.Database
	opts   Options
}

type Options struct {
	Notifications           string `yaml:"notifications"`
	NotificationPreferences string `yaml:"notification_preferences"`
	NotificationTemplates   string `yaml:"notification_templates"`
}

func InitNotificationRepository(redis0 *redis.Client, mongo0 *mongo.Database) NotificationRepositoryItf {
	return &notificationRepository{
		redis0: redis0,
		mongo0: mongo0,
	}
}

func replaceTemplate(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = replaceAll(result, placeholder, value)
	}

	return result
}

func replaceAll(s, old, new string) string {
	for i := 0; i < len(s); i++ {
		if len(s[i:]) >= len(old) && s[i:i+len(old)] == old {
			return s[:i] + new + replaceAll(s[i+len(old):], old, new)
		}
	}

	return s
}
