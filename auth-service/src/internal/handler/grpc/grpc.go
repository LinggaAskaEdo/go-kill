package grpc

import (
	"sync"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"github.com/rs/zerolog"
)

var onceGrpcHandler = &sync.Once{}

// type GrpcItf interface {
// 	CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error)
// 	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error)
// 	ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error)
// 	RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error)
// 	Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error)
// }

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
