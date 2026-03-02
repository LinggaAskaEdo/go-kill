package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository/auth"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceItf interface {
	CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error)
	ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error)
	RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error)
	Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error)
}

type authService struct {
	authRepository auth.AuthRepositoryItf
}

func InitAuthService(authRepository auth.AuthRepositoryItf) AuthServiceItf {
	return &authService{
		authRepository: authRepository,
	}
}

func generateTokenID() string {
	return fmt.Sprintf("token-%d", time.Now().UnixNano())
}

func generateRefreshToken() string {
	return fmt.Sprintf("refresh-%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

func hashToken(token string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(token), 10)
	return string(hash)
}
