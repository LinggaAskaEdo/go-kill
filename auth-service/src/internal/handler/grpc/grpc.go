package grpc

import (
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/rs/zerolog"
)

type GrpcItf interface {
	Serve() []string
}

type Grpc struct {
	authpb.AuthServiceServer
	log zerolog.Logger
	svc *service.Service
}

func InitGrpcHandler(log zerolog.Logger, svc *service.Service) *Grpc {
	return &Grpc{
		log: log,
		svc: svc,
	}
}
