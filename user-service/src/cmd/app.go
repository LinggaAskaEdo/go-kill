package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
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

	"golang.org/x/sync/errgroup"
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
	appSubComp := app.New(app.WithShutdownTimeout(15*time.Second), app.WithLogger(log))
	appMainComp := app.New(app.WithShutdownTimeout(15*time.Second), app.WithLogger(log))

	// Initialize database component
	dbComp0 := database.NewDatabaseComponent(log, cfg.Database["db-0"])
	if dbComp0 != nil {
		appSubComp.Add(dbComp0, 10*time.Second)
	}

	// Initialize query loader component
	queryComp := query.NewQueryComponent(log, cfg.Query)
	appSubComp.Add(queryComp, 10*time.Second)

	// Initialize mongo component
	mongoComp0 := mongo.NewMongoDBComponent(log, cfg.Mongo["mongo-0"])
	appSubComp.Add(mongoComp0, 10*time.Second)

	// Initialize gRPC client components
	authClientComp := grpcclient.NewGRPCClientComponent(log, cfg.GRPCClient["auth_service"])
	appSubComp.Add(authClientComp, 10*time.Second)

	// Initialze middleware
	mw := middleware.Init(log)

	// Initialize Gin engine
	gin := http.Init(log, mw, cfg.Http)

	// Stage 1: Start independent components (no dependencies)
	independent := []app.Component{dbComp0, queryComp, mongoComp0, authClientComp}

	// Create a shared context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Use an errgroup to manage their goroutines
	indepGroup, indepCtx := errgroup.WithContext(ctx)
	for _, comp := range independent {
		indepGroup.Go(func() error {
			return comp.Start(indepCtx)
		})
	}

	// Wait for each independent component to be ready (with timeout)
	for _, comp := range independent {
		select {
		case <-comp.Ready():
			log.Info().Str("component", fmt.Sprintf("%T", comp)).Msg("ready")
		case <-time.After(10 * time.Second):
			log.Fatal().Msgf("timeout waiting for component %T", comp)
		}
	}

	// Now build the service component (which depends on database, mongo, query, etc.)
	serviceComp := config.NewServiceComponent(log, dbComp0, queryComp, mongoComp0, authClientComp)
	appMainComp.Add(serviceComp, 10*time.Second)

	// Initialize scheduler component
	schedComp := scheduler.NewSchedulerComponent(log, func() ([]scheduler.Job, error) {
		select {
		case <-serviceComp.Ready():
			userGenJob := sched.NewUserGeneratorJob(log, serviceComp.Service().User, cfg.Scheduler["job-0"])
			return []scheduler.Job{userGenJob}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(10 * time.Second):
			return nil, fmt.Errorf("timeout waiting for scheduler")
		}
	})
	if schedComp != nil {
		appMainComp.Add(schedComp, 10*time.Second)
	}

	// Initialize gRPC server component
	grpcServerComp := grpcserver.NewGRPCServerComponent(log, cfg.GRPCServer, func(ctx context.Context, s *grpc.Server) error {
		select {
		case <-serviceComp.Ready():
			userpb.RegisterUserServiceServer(s, serviceComp.GrpcHandler())
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timeout waiting for gRPC Server")
		}
	})
	appMainComp.Add(grpcServerComp, 10*time.Second)

	// Initialize HTTP server component
	httpServerComp := server.NewHTTPServerComponent(log, cfg.Server, mw, gin, func(ctx context.Context, engine *server.Engine) error {
		select {
		case <-serviceComp.Ready():
			restHandler.InitRestHandler(engine, serviceComp.Service())
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timeout waiting for HTTP Server")
		}
	})
	appMainComp.Add(httpServerComp, 10*time.Second)

	// Run the app – now all components are added and will start in the order they were added.
	if err := appMainComp.Run(); err != nil {
		log.Fatal().Err(err).Msg("app failed")
	}

	if err := indepGroup.Wait(); err != nil && err != context.Canceled {
		log.Error().Err(err).Msg("Independent component error")
	}
}
