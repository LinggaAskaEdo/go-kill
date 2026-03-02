package grpc

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
)

func (g *Grpc) CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error) {
	return g.svc.Auth.CreateAuthUser(ctx, req)
}

func (g *Grpc) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return g.svc.Auth.Login(ctx, req)
}

func (g *Grpc) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	return g.svc.Auth.ValidateToken(ctx, req)
}

func (g *Grpc) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	return g.svc.Auth.RefreshToken(ctx, req)
}

func (g *Grpc) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	return g.svc.Auth.Logout(ctx, req)
}
