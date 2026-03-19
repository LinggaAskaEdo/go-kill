package repository

import (
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/repository/analytics"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	Analytics analytics.AnalyticsRepositoryItf
}

type Options struct {
	AnalyticsOpts analytics.Options `yaml:"analytics"`
}

func InitRepository(redis0 *redis.Client, mongo0 *mongo.Database, opts Options) *Repository {
	return &Repository{
		Analytics: analytics.InitAnalyticsRepository(
			redis0,
			mongo0,
			opts.AnalyticsOpts,
		),
	}
}
