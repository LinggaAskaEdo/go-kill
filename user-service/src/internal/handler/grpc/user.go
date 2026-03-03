package grpc

import (
	"context"

	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
)

func (g *Grpc) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	return g.svc.User.CreateUser(ctx, req)
}

func (g *Grpc) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return g.svc.User.GetUser(ctx, req)
}

func (g *Grpc) GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.GetAddressResponse, error) {
	return g.svc.User.GetAddress(ctx, req)
}

func (g *Grpc) LogActivity(ctx context.Context, req *userpb.LogActivityRequest) (*userpb.LogActivityResponse, error) {
	return g.svc.User.LogActivity(ctx, req)
}
