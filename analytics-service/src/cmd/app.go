package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/config"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/handler/pubsub"
	restHandler "github.com/linggaaskaedo/go-kill/analytics-service/src/internal/handler/rest"
	"github.com/linggaaskaedo/go-kill/common/app"
	"github.com/linggaaskaedo/go-kill/common/component/http"
	"github.com/linggaaskaedo/go-kill/common/component/kafkaconsumer"
	"github.com/linggaaskaedo/go-kill/common/component/kafkaproducer"
	"github.com/linggaaskaedo/go-kill/common/component/mongo"
	"github.com/linggaaskaedo/go-kill/common/component/redis"
	"github.com/linggaaskaedo/go-kill/common/component/server"
	"github.com/linggaaskaedo/go-kill/common/pkg/logger"
	"github.com/linggaaskaedo/go-kill/common/pkg/middleware"

	"golang.org/x/sync/errgroup"
)

var (
	minJitter int
	maxJitter int
)

// @title			Go-Kill x Analytic Service
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

	log.Info().Msg("Starting analytics service...")

	// Create application with options
	appSubComp := app.New(app.WithShutdownTimeout(15*time.Second), app.WithLogger(log))
	appMainComp := app.New(app.WithShutdownTimeout(15*time.Second), app.WithLogger(log))

	// Initialize redis component
	redisComp0 := redis.NewRedisComponent(log, cfg.Redis)
	appSubComp.Add(redisComp0, 10*time.Second)

	// Initialize database component
	mongoComp0 := mongo.NewMongoDBComponent(log, cfg.Mongo["mongo-0"])
	if mongoComp0 != nil {
		appSubComp.Add(mongoComp0, 10*time.Second)
	}

	// Initialize Kafka producer component
	producerComp := kafkaproducer.NewKafkaProducerComponent(log, cfg.KafkaProducer)
	appSubComp.Add(producerComp, 10*time.Second)

	// Stage 1: Start independent components (no dependencies)
	independent := []app.Component{redisComp0, mongoComp0, producerComp}

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

	select {
	case <-producerComp.Ready():
		log.Info().Str("component", fmt.Sprintf("%T", producerComp)).Msg("ready")
	case <-time.After(10 * time.Second):
		log.Fatal().Msgf("timeout waiting for component %T", producerComp)
	}

	// Initialize middleware
	mw := middleware.Init(log)

	// Initialize Gin engine
	gin := http.Init(log, mw, cfg.Http)

	// Now build the service component (which depends on database, mongo, query, etc.)
	serviceComp := config.NewServiceComponent(log, redisComp0, mongoComp0, cfg.Repository)
	appMainComp.Add(serviceComp, 10*time.Second)

	// Build Kafka consumer with producer for DLQ
	consumerHandler := pubsub.NewConsumerGroupHandler(log, serviceComp.Service(), producerComp.Producer())
	consumerComp := kafkaconsumer.NewKafkaConsumerComponent(log, cfg.KafkaConsumer, consumerHandler)
	appMainComp.Add(consumerComp, 10*time.Second)

	// Build HTTP server (depends on service)
	httpServerComp := server.NewHTTPServerComponent(log, cfg.Server, mw, gin, func(ctx context.Context, engine *server.Engine) error {
		select {
		case <-serviceComp.Ready():
			restHandler.InitRestHandler(engine, serviceComp.Service(), serviceComp.Redis(), serviceComp.Mongo())
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
