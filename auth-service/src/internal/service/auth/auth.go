package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository/auth"
)

type AuthServiceItf interface {
	CreateAuthUser(ctx context.Context, req *dto.CreateAuthUserRequest) (*dto.CreateAuthUserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	ValidateToken(ctx context.Context, req *dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error)
}

type authService struct {
	authRepository auth.AuthRepositoryItf
	authOptions    Options
}

type Options struct {
	JwtSecret string `yaml:"jwt_secret"`
}

func InitAuthService(authRepository auth.AuthRepositoryItf, authOptions Options) AuthServiceItf {
	return &authService{
		authRepository: authRepository,
		authOptions:    authOptions,
	}
}

func generateTokenID() string {
	return fmt.Sprintf("token-%d", time.Now().UnixNano())
}

func generateRefreshToken() string {
	return fmt.Sprintf("refresh-%d-%d", time.Now().UnixNano(), time.Now().Unix())
}
