package grpc

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
)

func (g *Grpc) CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error) {
	authResp, err := g.authClient.CreateAuthUser(ctx, &authpb.CreateAuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &authpb.CreateAuthUserResponse{
		AuthId: authResp.AuthId,
	}, nil
}
