package config

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/grpcclient"
	"github.com/linggaaskaedo/go-kill/common/component/kafkaproducer"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	grpcHandler "github.com/linggaaskaedo/go-kill/order-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service"

	"github.com/rs/zerolog"
)

type ServiceComponent struct {
	log               zerolog.Logger
	dbComp0           *database.DatabaseComponent
	queryComp         *query.QueryComponent
	userClientComp    *grpcclient.GRPCClientComponent
	productClientComp *grpcclient.GRPCClientComponent
	kafkaProducerComp *kafkaproducer.KafkaProducerComponent

	repo        *repository.Repository
	service     *service.Service
	grpcHandler *grpcHandler.Grpc
	ready       chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	dbComp0 *database.DatabaseComponent,
	queryComp *query.QueryComponent,
	userClientComp *grpcclient.GRPCClientComponent,
	productClientComp *grpcclient.GRPCClientComponent,
	kafkaProducerComp *kafkaproducer.KafkaProducerComponent,
) *ServiceComponent {
	return &ServiceComponent{
		log:               log,
		dbComp0:           dbComp0,
		queryComp:         queryComp,
		userClientComp:    userClientComp,
		productClientComp: productClientComp,
		kafkaProducerComp: kafkaProducerComp,
		ready:             make(chan struct{}),
	}
}

func (s *ServiceComponent) Start(ctx context.Context) error {
	s.repo = repository.InitRepository(s.dbComp0.Client(), s.queryComp, s.productClientComp.Conn())
	s.service = service.InitService(s.repo, s.userClientComp.Conn(), s.productClientComp.Conn(), s.kafkaProducerComp)
	s.grpcHandler = grpcHandler.InitGrpcHandler(s.log, s.service)

	close(s.ready) // signal that service is ready
	s.log.Debug().Msg("Service component started")
	<-ctx.Done()

	return nil
}

func (s *ServiceComponent) Stop(ctx context.Context) error {
	s.log.Info().Msg("Service component stopping")

	if err := s.userClientComp.Stop(ctx); err != nil {
		s.log.Error().Err(err).Msg("failed to stop user client")
	}
	if err := s.productClientComp.Stop(ctx); err != nil {
		s.log.Error().Err(err).Msg("failed to stop product client")
	}
	if err := s.kafkaProducerComp.Stop(ctx); err != nil {
		s.log.Error().Err(err).Msg("failed to stop kafka producer")
	}
	if err := s.dbComp0.Stop(ctx); err != nil {
		s.log.Error().Err(err).Msg("failed to stop database")
	}

	s.log.Info().Msg("Service component stopped")
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
