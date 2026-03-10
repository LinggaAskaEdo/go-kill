package analytics

import (
	"context"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AnalyticsRepositoryItf interface {
	UpdateOrderAnalytics(ctx context.Context, event dto.OrderEvent) error
	UpdateProductAnalytics(ctx context.Context, event dto.OrderEvent) error
	UpdateCancellationMetrics(ctx context.Context, event dto.OrderEvent) error
}

type analyticsRepository struct {
	redis0           *redis.Client
	mongo0           *mongo.Database
	analyticsOptions Options
}

type Options struct {
	OrderCollection   string `yaml:"order_colllection"`
	ProductCollection string `yaml:"prodcut_collection"`
}

func InitAnalyticsRepository(redis0 *redis.Client, mongo0 *mongo.Database) AnalyticsRepositoryItf {
	return &analyticsRepository{
		redis0: redis0,
		mongo0: mongo0,
	}
}
