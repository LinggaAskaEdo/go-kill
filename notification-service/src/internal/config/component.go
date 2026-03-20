package config

import (
	"context"
	"fmt"

	"github.com/linggaaskaedo/go-kill/common/component/mongo"
	"github.com/linggaaskaedo/go-kill/common/component/redis"
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/service"

	"github.com/rs/zerolog"
)

type ServiceComponent struct {
	log        zerolog.Logger
	redisComp0 *redis.RedisComponent
	mongoComp0 *mongo.MongoDBComponent
	repoOpts   repository.Options

	repo    *repository.Repository
	service *service.Service
	ready   chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	redisComp0 *redis.RedisComponent,
	mongoComp0 *mongo.MongoDBComponent,
	repoOpts repository.Options,
) *ServiceComponent {
	return &ServiceComponent{
		log:        log,
		redisComp0: redisComp0,
		mongoComp0: mongoComp0,
		repoOpts:   repoOpts,
		ready:      make(chan struct{}),
	}
}

func (s *ServiceComponent) Start(ctx context.Context) error {
	if s.mongoComp0 == nil {
		return fmt.Errorf("mongo component is nil")
	}

	s.repo = repository.InitRepository(s.redisComp0.Client(), s.mongoComp0.Database(), s.repoOpts)
	if s.repo == nil {
		return fmt.Errorf("failed to init repository")
	}

	s.service = service.InitService(s.repo)
	if s.service == nil {
		return fmt.Errorf("failed to init service")
	}

	close(s.ready)
	s.log.Debug().Msg("Service component started")
	<-ctx.Done()

	return nil
}

func (s *ServiceComponent) Stop(ctx context.Context) error {
	s.log.Debug().Msg("Service component stopping")

	if s.repo != nil {
		if s.repo.Notification != nil {
			s.log.Debug().Msg("Service component stopped")
		}
	}

	return nil
}

func (s *ServiceComponent) Service() *service.Service {
	if s.service == nil {
		panic("ServiceComponent.Service() called before Start()")
	}

	return s.service
}

func (s *ServiceComponent) Ready() <-chan struct{} {
	return s.ready
}
