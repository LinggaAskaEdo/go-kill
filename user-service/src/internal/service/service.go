package service

import (
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service/user"
)

type Service struct {
	User user.UserServiceItf
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		User: user.InitUserService(
			repository.User,
		),
	}
}
