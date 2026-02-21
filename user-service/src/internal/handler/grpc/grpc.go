package grpc

import (
	authpb "github.com/linggaaskaedo/go-kill/user-service/src/api/proto"
	userpb "github.com/linggaaskaedo/go-kill/user-service/src/api/proto"
)

type Grpc struct {
	userpb.UnimplementedUserServiceServer
	authClient authpb.AuthServiceClient
}

func IniGrpcHandler(authClient authpb.AuthServiceClient) *Grpc {
	return &Grpc{
		authClient: authClient,
	}
}
