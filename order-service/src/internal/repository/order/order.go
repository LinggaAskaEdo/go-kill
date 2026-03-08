package order

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/query"
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type OrderRepositoryItf interface {
	StoreOrder(ctx context.Context, productDetails []*dto.ProductDetails, createOrders *dto.CreateOrderRequest, totalAmount float64) (*string, *string, error)
	GetOrder(ctx context.Context, reqData *dto.GetOrderRequest) (*entity.Order, error)
	ListOrders(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, int32, error)
	CancelOrder(ctx context.Context, reqData *dto.CancelOrderRequest) error
}

type orderRepository struct {
	db0           *sqlx.DB
	queryLoader   *query.QueryComponent
	productClient productpb.ProductServiceClient
}

func InitOrderRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, productClientConn *grpc.ClientConn) OrderRepositoryItf {
	return &orderRepository{
		db0:           db0,
		queryLoader:   queryLoader,
		productClient: productpb.NewProductServiceClient(productClientConn),
	}
}
