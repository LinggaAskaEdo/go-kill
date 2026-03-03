package service

import (
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service/user"

	"google.golang.org/grpc"
)

type Service struct {
	User user.UserServiceItf
}

func InitService(clientConn *grpc.ClientConn, repository *repository.Repository) *Service {
	return &Service{
		User: user.InitUserService(
			clientConn,
			repository.User,
		),
	}
}
