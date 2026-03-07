package service

import (
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service/order"
)

type Service struct {
	Order order.OrderServiceItf
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		Order: order.InitOrderService(
			repository.Order,
		),
	}
}
