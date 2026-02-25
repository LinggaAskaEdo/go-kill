package grpc

import (
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
)

type Grpc struct {
	userpb.UnimplementedUserServiceServer
	authClient authpb.AuthServiceClient
}

func InitGrpcHandler(authClient authpb.AuthServiceClient) *Grpc {
	return &Grpc{
		authClient: authClient,
	}
}
