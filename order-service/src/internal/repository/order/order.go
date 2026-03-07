package order

import (
	"github.com/linggaaskaedo/go-kill/common/component/query"

	"github.com/jmoiron/sqlx"
)

type OrderRepositoryItf interface {
}

type orderRepository struct {
	db0         *sqlx.DB
	queryLoader *query.QueryComponent
}

func InitOrderRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent) OrderRepositoryItf {
	return &orderRepository{
		db0:         db0,
		queryLoader: queryLoader,
	}
}
