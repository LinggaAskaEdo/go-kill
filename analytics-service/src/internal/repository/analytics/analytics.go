package analytics

import (
	"context"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AnalyticsRepositoryItf interface {
	UpdateOrderAnalytics(ctx context.Context, event dto.OrderEvent) error
	UpdateProductAnalytics(ctx context.Context, event dto.OrderEvent) error
	UpdateCancellationMetrics(ctx context.Context, event dto.OrderEvent) error
	EnsureIndexes(ctx context.Context) error
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

func InitAnalyticsRepository(redis0 *redis.Client, mongo0 *mongo.Database, opts Options) AnalyticsRepositoryItf {
	return &analyticsRepository{
		redis0:           redis0,
		mongo0:           mongo0,
		analyticsOptions: opts,
	}
}

func (r *analyticsRepository) EnsureIndexes(ctx context.Context) error {
	orderColl := r.mongo0.Collection(r.analyticsOptions.OrderCollection)
	productColl := r.mongo0.Collection(r.analyticsOptions.ProductCollection)

	orderIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "date", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	productIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}, {Key: "date", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := orderColl.Indexes().CreateMany(ctx, orderIndexes)
	if err != nil {
		return err
	}

	_, err = productColl.Indexes().CreateMany(ctx, productIndexes)
	if err != nil {
		return err
	}

	return nil
}
