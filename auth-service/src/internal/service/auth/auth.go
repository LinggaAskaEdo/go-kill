package auth

import (
	"context"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository/auth"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
)

type AuthServiceItf interface {
	CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error)
}

type authService struct {
	authRepository auth.AuthRepositoryItf
}

func InitAuthService(authRepository auth.AuthRepositoryItf) AuthServiceItf {
	return &authService{
		authRepository: authRepository,
	}
}
