package service

import (
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service/auth"
)

type Service struct {
	Auth auth.AuthServiceItf
}

type Options struct {
	AuthOpts auth.Options `yaml:"auth"`
}

func InitService(repository *repository.Repository, opts Options) *Service {
	return &Service{
		Auth: auth.InitAuthService(
			repository.Auth,
			opts.AuthOpts,
		),
	}
}
