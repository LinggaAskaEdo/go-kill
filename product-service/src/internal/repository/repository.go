package repository

import (
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/repository/product"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Product product.ProductRepositoryItf
}

func InitRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, redis0 *redis.Client) *Repository {
	return &Repository{
		Product: product.InitProductRepository(
			db0,
			queryLoader,
			redis0,
		),
	}
}
