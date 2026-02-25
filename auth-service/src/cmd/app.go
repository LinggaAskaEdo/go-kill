package main

import (
	"flag"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/config"
	grpcHandler "github.com/linggaaskaedo/go-kill/auth-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/common/app"
	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/grpcserver"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/common/component/redis"
	"github.com/linggaaskaedo/go-kill/common/pkg/logger"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"google.golang.org/grpc"
)

var (
	minJitter int
	maxJitter int
)

// @title			Go-Kill x Auth Service
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

	// Initialize redis component
	redisComp0 := redis.NewRedisComponent(log, cfg.Redis, "apps")
	application.Add(redisComp0, 10*time.Second)

	// Initialize database component
	dbComp0 := database.NewDatabaseComponent(log, cfg.Database["db-0"])
	if dbComp0 != nil {
		application.Add(dbComp0, 10*time.Second)
	}

	// Initialize query loader component
	queryComp := query.NewQueryComponent(log, cfg.Query)
	application.Add(queryComp, 10*time.Second)

	// Initialize service component
	serviceComp := config.NewServiceComponent(log, dbComp0, queryComp, redisComp0)
	application.Add(serviceComp, 10*time.Second)

	// Initialize gRPC server component
	grpcServerComp := grpcserver.NewGRPCServerComponent(log, cfg.GRPCServer, func(s *grpc.Server) {
		log.Info().Msg("gRPC server registrar: waiting for service component...")
		<-serviceComp.Ready()
		log.Info().Msg("gRPC server registrar: service ready, creating handler")
		grpcHandler := grpcHandler.InitGrpcHandler(log, serviceComp.Service())
		authpb.RegisterAuthServiceServer(s, grpcHandler)
		log.Info().Msg("gRPC server registrar: handler registered")
	})
	application.Add(grpcServerComp, 10*time.Second)

	// Run Server
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}
}
