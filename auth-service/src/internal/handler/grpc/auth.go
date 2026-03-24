package grpc

import (
	"context"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
)

func (g *Grpc) CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error) {
	dtoReq := &dto.CreateAuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := g.svc.Auth.CreateAuthUser(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.CreateAuthUserResponse{
		Success: resp.Success,
		AuthId:  resp.AuthId,
	}, nil
}

func (g *Grpc) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	dtoReq := &dto.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		IpAddress: req.IpAddress,
		UserAgent: req.UserAgent,
	}

	resp, err := g.svc.Auth.Login(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Success:      resp.Success,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (g *Grpc) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	dtoReq := &dto.ValidateTokenRequest{
		Token: req.Token,
	}

	resp, err := g.svc.Auth.ValidateToken(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.ValidateTokenResponse{
		Valid:  resp.Valid,
		UserId: resp.UserId,
		Email:  resp.Email,
	}, nil
}

func (g *Grpc) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	dtoReq := &dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	resp, err := g.svc.Auth.RefreshToken(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshTokenResponse{
		Success:      resp.Success,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (g *Grpc) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	dtoReq := &dto.LogoutRequest{
		Token:  req.Token,
		UserId: req.UserId,
	}

	resp, err := g.svc.Auth.Logout(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.LogoutResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}
