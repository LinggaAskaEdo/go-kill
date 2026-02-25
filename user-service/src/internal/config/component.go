package config

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/grpcclient"
	"github.com/linggaaskaedo/go-kill/common/component/mongo"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	grpcHandler "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"

	"github.com/rs/zerolog"
)

type ServiceComponent struct {
	log            zerolog.Logger
	dbComp0        *database.DatabaseComponent
	queryComp      *query.QueryComponent
	mongoComp0     *mongo.MongoDBComponent
	authClientComp *grpcclient.GRPCClientComponent

	repo        *repository.Repository
	service     *service.Service
	grpcHandler *grpcHandler.Grpc
	ready       chan struct{}
}

func NewServiceComponent(
	log zerolog.Logger,
	dbComp0 *database.DatabaseComponent,
	queryComp *query.QueryComponent,
	mongoComp0 *mongo.MongoDBComponent,
	authClientComp *grpcclient.GRPCClientComponent,
) *ServiceComponent {
	return &ServiceComponent{
		log:            log,
		dbComp0:        dbComp0,
		queryComp:      queryComp,
		mongoComp0:     mongoComp0,
		authClientComp: authClientComp,
		ready:          make(chan struct{}),
	}
}

func (s *ServiceComponent) Start(ctx context.Context) error {
	select {
	case <-s.dbComp0.Ready():
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case <-s.queryComp.Ready():
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case <-s.mongoComp0.Ready():
	case <-ctx.Done():
		return ctx.Err()
	}

	db := s.dbComp0.Client()
	ql := s.queryComp
	mongoDB := s.mongoComp0.Database()
	s.repo = repository.InitRepository(db, ql, mongoDB)
	s.service = service.InitService(s.repo)

	// Wait for auth client (if enabled) – but it may fail; respect context
	if s.authClientComp != nil {
		select {
		case <-s.authClientComp.Ready():
			authClient := authpb.NewAuthServiceClient(s.authClientComp.Conn())
			s.grpcHandler = grpcHandler.InitGrpcHandler(authClient)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

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
