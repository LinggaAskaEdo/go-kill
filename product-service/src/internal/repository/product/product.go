package product

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ProductRepositoryItf interface {
	CreateProduct(ctx context.Context, product *entity.Product, qty int, rsv int) (*entity.Product, error)
	GetListProduct(ctx context.Context) ([]*entity.Product, error)
	GetProduct(ctx context.Context, productID string) (*entity.Product, error)
	ListCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategoriesByProduct(ctx context.Context, productID string) ([]*entity.Category, error)
	GetProductsByCategory(ctx context.Context, categoryID string) ([]*entity.Product, error)
}

type productRepository struct {
	db0         *sqlx.DB
	queryLoader *query.QueryComponent
	redis0      *redis.Client
}

func InitProductRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, redis0 *redis.Client) ProductRepositoryItf {
	return &productRepository{
		db0:         db0,
		queryLoader: queryLoader,
		redis0:      redis0,
	}
}
