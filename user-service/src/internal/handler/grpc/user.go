package grpc

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/user-service/src/api/proto"
	userpb "github.com/linggaaskaedo/go-kill/user-service/src/api/proto"
)

func (g *Grpc) CreateAuthUser(ctx context.Context, req *userpb.CreateAuthUserRequest) (*userpb.CreateAuthUserResponse, error) {
	// Use the auth client
	authResp, err := g.authClient.CreateAuthUser(ctx, &authpb.CreateAuthUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.CreateAuthUserResponse{AuthId: authResp.AuthId}, nil
}
