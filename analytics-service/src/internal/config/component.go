package config

import (
	"context"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service"
	mongocomponent "github.com/linggaaskaedo/go-kill/common/component/mongo"
	rediscomponent "github.com/linggaaskaedo/go-kill/common/component/redis"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rs/zerolog"
)

type ServiceComponent struct {
	log        zerolog.Logger
	redisComp0 *rediscomponent.RedisComponent
	mongoComp0 *mongocomponent.MongoDBComponent
	repoOpts   repository.Options

	repo    *repository.Repository
	service *service.Service
	ready   chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	redisComp0 *rediscomponent.RedisComponent,
	mongoComp0 *mongocomponent.MongoDBComponent,
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
	s.repo = repository.InitRepository(s.redisComp0.Client(), s.mongoComp0.Database(), s.repoOpts)

	if err := s.repo.Analytics.EnsureIndexes(ctx); err != nil {
		s.log.Warn().Err(err).Msg("Failed to ensure indexes, continuing anyway")
	}

	s.service = service.InitService(s.repo)

	close(s.ready)
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

func (s *ServiceComponent) Redis() *goredis.Client {
	return s.redisComp0.Client()
}

func (s *ServiceComponent) Mongo() *mongo.Database {
	return s.mongoComp0.Database()
}
