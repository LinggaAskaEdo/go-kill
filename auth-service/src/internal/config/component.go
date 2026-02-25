package config

import (
	"context"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/common/component/redis"

	"github.com/rs/zerolog"
)

type ServiceComponent struct {
	log        zerolog.Logger
	dbComp0    *database.DatabaseComponent
	queryComp  *query.QueryComponent
	redisComp0 *redis.RedisComponent
	repo       *repository.Repository
	service    *service.Service
	ready      chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	dbComp0 *database.DatabaseComponent,
	queryComp *query.QueryComponent,
	redisComp0 *redis.RedisComponent,
) *ServiceComponent {
	return &ServiceComponent{
		log:        log,
		dbComp0:    dbComp0,
		queryComp:  queryComp,
		redisComp0: redisComp0,
		ready:      make(chan struct{}),
	}
}

func (s *ServiceComponent) Start(ctx context.Context) error {
	// Wait for all dependencies to be ready
	<-s.dbComp0.Ready()
	<-s.queryComp.Ready()
	<-s.redisComp0.Ready()

	// Now construct repository and service
	db := s.dbComp0.Client()
	ql := s.queryComp
	rd := s.redisComp0.Client()

	s.repo = repository.InitRepository(db, ql, rd)
	s.service = service.InitService(s.repo)

	close(s.ready) // signal that service is ready
	s.log.Debug().Msg("Service component started")
	<-ctx.Done()

	return nil
}

func (s *ServiceComponent) Stop(ctx context.Context) error {
	s.log.Debug().Msg("Service component stopped")
	return nil
}

func (s *ServiceComponent) Service() *service.Service {
	return s.service
}

func (s *ServiceComponent) Ready() <-chan struct{} {
	return s.ready
}
