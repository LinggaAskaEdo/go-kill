package grpc

import (
	"sync"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/rs/zerolog"
)

var onceGrpcHandler = &sync.Once{}

type Grpc struct {
	authpb.AuthServiceServer
	log zerolog.Logger
	svc *service.Service
}

func InitGrpcHandler(log zerolog.Logger, svc *service.Service) *Grpc {
	var grpcHandler *Grpc

	onceGrpcHandler.Do(func() {
		grpcHandler = &Grpc{
			log: log,
			svc: svc,
		}
	})

	return grpcHandler
}
