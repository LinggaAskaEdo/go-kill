package config

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	grpcHandler "github.com/linggaaskaedo/go-kill/order-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service"

	"github.com/rs/zerolog"
)

type ServiceComponent struct {
	log       zerolog.Logger
	dbComp0   *database.DatabaseComponent
	queryComp *query.QueryComponent

	repo        *repository.Repository
	service     *service.Service
	grpcHandler *grpcHandler.Grpc
	ready       chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	dbComp0 *database.DatabaseComponent,
	queryComp *query.QueryComponent,
) *ServiceComponent {
	return &ServiceComponent{
		log:       log,
		dbComp0:   dbComp0,
		queryComp: queryComp,
		ready:     make(chan struct{}),
	}
}

func (s *ServiceComponent) Start(ctx context.Context) error {
	s.repo = repository.InitRepository(s.dbComp0.Client(), s.queryComp)
	s.service = service.InitService(s.repo)
	s.grpcHandler = grpcHandler.InitGrpcHandler(s.log, s.service)

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

func (s *ServiceComponent) GrpcHandler() *grpcHandler.Grpc {
	return s.grpcHandler
}

func (s *ServiceComponent) Ready() <-chan struct{} {
	return s.ready
}
