package main

import (
	"flag"
	"time"

	"github.com/linggaaskaedo/go-kill/common/app"
	"github.com/linggaaskaedo/go-kill/common/database"
	"github.com/linggaaskaedo/go-kill/common/grpcclient"
	"github.com/linggaaskaedo/go-kill/common/grpcserver"
	"github.com/linggaaskaedo/go-kill/common/http"
	"github.com/linggaaskaedo/go-kill/common/logger"
	"github.com/linggaaskaedo/go-kill/common/middleware"
	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/linggaaskaedo/go-kill/common/redis"
	"github.com/linggaaskaedo/go-kill/common/scheduler"
	"github.com/linggaaskaedo/go-kill/common/server"
	authpb "github.com/linggaaskaedo/go-kill/user-service/src/api/proto"
	userpb "github.com/linggaaskaedo/go-kill/user-service/src/api/proto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/config"
	grpcHandler "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	restHandler "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/rest"
	sched "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/scheduler"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"

	"google.golang.org/grpc"
)

var (
	minJitter int
	maxJitter int
)

// @title			Go-Kill
// @version		1.0
// @description	Microservices Architecture with Go
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host			localhost:8080
// @schemes		http https
func main() {
	flag.IntVar(&minJitter, "minSleep", DefaultMinJitter, "min. sleep duration during app initialization")
	flag.IntVar(&maxJitter, "maxSleep", DefaultMaxJitter, "max. sleep duration during app initialization")
	flag.Parse()

	// Add sleep with Jitter to drag the the initialization time among instances
	sleepWithJitter(minJitter, maxJitter)

	// Load config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	// Initialize logger
	log := logger.Init(cfg.Logger)

	log.Info().Msg("Starting user service...")

	// Create application with options
	application := app.New(app.WithShutdownTimeout(15*time.Second), app.WithLogger(log))

	// Initialize Redis Component
	redisComp := redis.NewRedisComponent(log, cfg.Redis, "apps")
	application.Add(redisComp, 10*time.Second)

	// Initialize database components
	dbComp0 := database.NewDatabaseComponent(log, cfg.Database["db-0"])
	if dbComp0 != nil {
		application.Add(dbComp0, 10*time.Second)
	}

	// Initialize query loader component
	queryComp := query.NewQueryComponent(log, cfg.Query)
	application.Add(queryComp, 10*time.Second)

	// Initialize dependencies
	repository := repository.InitRepository(dbComp0.Client(), queryComp)
	service := service.InitService(repository)

	// Initialize scheduler component
	userGenJob := sched.NewUserGeneratorJob(log, service.User, cfg.Scheduler.SchedulerJobs.UserGeneratorJob)
	jobs := []scheduler.Job{userGenJob}

	schedComp := scheduler.NewSchedulerComponent(log, cfg.Scheduler, jobs)
	if schedComp != nil {
		application.Add(schedComp, 10*time.Second)
	}

	// Initialze middleware
	mw := middleware.Init(log)

	// Initialize Gin engine
	gin := http.Init(log, mw, cfg.Http)

	// Initialize gRPC client components
	authClientComp := grpcclient.NewGRPCClientComponent(log, cfg.GRPCClient["auth_service"])
	application.Add(authClientComp, 10*time.Second)

	grpcHandler := grpcHandler.IniGrpcHandler(authpb.NewAuthServiceClient(authClientComp.Conn()))

	userServerComp := grpcserver.NewGRPCServerComponent(log, cfg.GRPCServer, func(s *grpc.Server) {
		userpb.RegisterUserServiceServer(s, grpcHandler)
	})
	application.Add(userServerComp, 10*time.Second)

	// Initialize REST handlers
	restHandler.InitRestHandler(gin, service, grpcHandler)

	// Initialize HTTP server component
	serverComp := server.NewHTTPServerComponent(log, cfg.Server, mw, gin)
	application.Add(serverComp, 10*time.Second)

	// Connect to Auth Service
	// const authServiceAddr = "localhost:50051"
	// authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// authClient := authpb.NewAuthServiceClient(authConn)
	// authClient.CreateAuthUser(context.Background(), &authpb.CreateAuthUserRequest{Email: "test@example.com", Password: "password"})

	// go func() {
	// 	lis, err := net.Listen("tcp", ":8082")
	// 	if err != nil {
	// 		log.Fatal().Err(err).Msg("Failed to listen on port 8082")
	// 	}

	// 	s := grpc.NewServer()
	// 	userpb.RegisterUserServiceServer(s, userpb.UnimplementedUserServiceServer{})

	// 	log.Println("User Service gRPC listening on :8082")
	// 	if err := s.Serve(lis); err != nil {
	// 		log.Fatal().Err(err).Msg("Failed to serve gRPC")
	// 	}
	// }()

	// Run Server
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}
}
