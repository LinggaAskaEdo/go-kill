package order

import "github.com/linggaaskaedo/go-kill/order-service/src/internal/repository/order"

type OrderServiceItf interface {
}

type orderService struct {
	orderRepository order.OrderRepositoryItf
}

func InitOrderService(orderRepository order.OrderRepositoryItf) OrderServiceItf {
	return &orderService{
		orderRepository: orderRepository,
	}
}
