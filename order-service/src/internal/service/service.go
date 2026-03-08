package service

import (
	"github.com/linggaaskaedo/go-kill/common/component/kafkaproducer"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service/order"

	"google.golang.org/grpc"
)

type Service struct {
	Order order.OrderServiceItf
}

type Options struct {
	OrderOpts order.Options `yaml:"order"`
}

func InitService(repository *repository.Repository, userClientConn *grpc.ClientConn, productClientConn *grpc.ClientConn, kafkaProducer *kafkaproducer.KafkaProducerComponent) *Service {
	return &Service{
		Order: order.InitOrderService(
			repository.Order,
			userClientConn,
			productClientConn,
			kafkaProducer,
		),
	}
}
