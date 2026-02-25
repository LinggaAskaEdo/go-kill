package main

import (
	"flag"
	"time"

	"github.com/linggaaskaedo/go-kill/common/app"
	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/grpcclient"
	"github.com/linggaaskaedo/go-kill/common/component/grpcserver"
	"github.com/linggaaskaedo/go-kill/common/component/http"
	"github.com/linggaaskaedo/go-kill/common/component/mongo"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/common/component/scheduler"
	"github.com/linggaaskaedo/go-kill/common/component/server"
	"github.com/linggaaskaedo/go-kill/common/pkg/logger"
	"github.com/linggaaskaedo/go-kill/common/pkg/middleware"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/config"
	restHandler "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/rest"
	sched "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/scheduler"

	"google.golang.org/grpc"
)

var (
	minJitter int
	maxJitter int
)

// @title			Go-Kill x User Service
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

	// Initialize database component
	dbComp0 := database.NewDatabaseComponent(log, cfg.Database["db-0"])
	if dbComp0 != nil {
		application.Add(dbComp0, 10*time.Second)
	}

	// Initialize query loader component
	queryComp := query.NewQueryComponent(log, cfg.Query)
	application.Add(queryComp, 10*time.Second)

	// Initialize mongo component
	mongoComp0 := mongo.NewMongoDBComponent(log, cfg.Mongo["mongo-0"])
	application.Add(mongoComp0, 10*time.Second)

	// Initialize gRPC client components
	authClientComp := grpcclient.NewGRPCClientComponent(log, cfg.GRPCClient["auth_service"])
	application.Add(authClientComp, 10*time.Second)

	// Initialize service component
	serviceComp := config.NewServiceComponent(log, dbComp0, queryComp, mongoComp0, authClientComp)
	application.Add(serviceComp, 10*time.Second)

	// Initialize scheduler component
	schedComp := scheduler.NewSchedulerComponent(log, cfg.Scheduler, func() ([]scheduler.Job, error) {
		<-serviceComp.Ready()
		userGenJob := sched.NewUserGeneratorJob(log, serviceComp.Service().User, cfg.Scheduler.SchedulerJobs.UserGeneratorJob)

		return []scheduler.Job{userGenJob}, nil
	})
	if schedComp != nil {
		application.Add(schedComp, 10*time.Second)
	}

	// Initialze middleware
	mw := middleware.Init(log)

	// Initialize Gin engine
	gin := http.Init(log, mw, cfg.Http)

	// Initialize gRPC server component
	grpcServerComp := grpcserver.NewGRPCServerComponent(log, cfg.GRPCServer, func(s *grpc.Server) {
		select {
		case <-serviceComp.Ready():
			userpb.RegisterUserServiceServer(s, serviceComp.GrpcHandler())
		case <-time.After(10 * time.Second):
			log.Error().Msg("Timed out waiting for service component; gRPC server will have no services")
		}
	})
	application.Add(grpcServerComp, 10*time.Second)

	// Initialize HTTP server component
	httpServerComp := server.NewHTTPServerComponent(log, cfg.Server, mw, gin, func(engine *server.Engine) {
		<-serviceComp.Ready()
		restHandler.InitRestHandler(gin, serviceComp.Service(), serviceComp.GrpcHandler())
	})
	application.Add(httpServerComp, 10*time.Second)

	// Run Server
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}
}
