package config

import (
	"context"

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

	repo    *repository.Repository
	service *service.Service
	ready   chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	redisComp0 *redis.RedisComponent,
	mongoComp0 *mongo.MongoDBComponent,
) *ServiceComponent {
	return &ServiceComponent{
		log:        log,
		redisComp0: redisComp0,
		mongoComp0: mongoComp0,
		ready:      make(chan struct{}),
	}
}

func (s *ServiceComponent) Start(ctx context.Context) error {
	s.repo = repository.InitRepository(s.redisComp0.Client(), s.mongoComp0.Database())
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
