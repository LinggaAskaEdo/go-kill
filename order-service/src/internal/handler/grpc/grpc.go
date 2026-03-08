package grpc

import (
	orderpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/order"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service"

	"github.com/rs/zerolog"
)

type Grpc struct {
	orderpb.OrderServiceServer
	log zerolog.Logger
	svc *service.Service
}

func InitGrpcHandler(log zerolog.Logger, svc *service.Service) *Grpc {
	return &Grpc{
		log: log,
		svc: svc,
	}
}
