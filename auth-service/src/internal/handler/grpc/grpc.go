package grpc

import (
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/rs/zerolog"
)

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

func (g *Grpc) Serve() []string {
	return []string{
		"/auth.AuthService/CreateAuthUser",
		"/auth.AuthService/Login",
		"/auth.AuthService/ValidateToken",
		"/auth.AuthService/RefreshToken",
		"/auth.AuthService/Logout",
	}
}
