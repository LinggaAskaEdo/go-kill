package repository

import (
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository/order"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type Repository struct {
	Order order.OrderRepositoryItf
}

func InitRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, productClientConn *grpc.ClientConn) *Repository {
	return &Repository{
		Order: order.InitOrderRepository(
			db0,
			queryLoader,
			productClientConn,
		),
	}
}
