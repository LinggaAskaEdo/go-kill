package service

import (
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service/auth"
)

type Service struct {
	Auth auth.AuthServiceItf
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		Auth: auth.InitAuthService(
			repository.Auth,
		),
	}
}
