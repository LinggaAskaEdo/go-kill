package service

import (
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service/analytics"
)

type Service struct {
	Analytics analytics.AnalyticsServiceItf
}

type Options struct {
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		Analytics: analytics.InitAnalyticsService(
			repository.Analytics,
		),
	}
}
