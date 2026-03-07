package order

import (
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository/order"

	"google.golang.org/grpc"
)

type OrderServiceItf interface {
}

type orderService struct {
	userClient      userpb.UserServiceClient
	productClient   productpb.ProductServiceClient
	orderRepository order.OrderRepositoryItf
}

func InitOrderService(userClientConn *grpc.ClientConn, productClientConn *grpc.ClientConn, orderRepository order.OrderRepositoryItf) OrderServiceItf {
	return &orderService{
		userClient:      userpb.NewUserServiceClient(userClientConn),
		productClient:   productpb.NewProductServiceClient(productClientConn),
		orderRepository: orderRepository,
	}
}
