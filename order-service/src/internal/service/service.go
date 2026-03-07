package service

import (
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service/order"

	"google.golang.org/grpc"
)

type Service struct {
	Order order.OrderServiceItf
}

func InitService(userClientConn *grpc.ClientConn, productClientConn *grpc.ClientConn, repository *repository.Repository) *Service {
	return &Service{
		Order: order.InitOrderService(
			userClientConn,
			productClientConn,
			repository.Order,
		),
	}
}
