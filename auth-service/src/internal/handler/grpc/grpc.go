package grpc

import (
	"sync"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/rs/zerolog"
)

var onceGrpcHandler = &sync.Once{}

type GrpcItf interface {
	Serve() []string
}

type Grpc struct {
	authpb.AuthServiceServer
	log zerolog.Logger
	svc *service.Service
}

func InitGrpcHandler(log zerolog.Logger, svc *service.Service) *Grpc {
	var g *Grpc

	onceGrpcHandler.Do(func() {
		g = &Grpc{
			log: log,
			svc: svc,
		}

		_ = g.Serve()
	})

	return g
}
