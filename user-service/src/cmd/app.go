package main

import (
	"flag"
	"time"

	"github.com/linggaaskaedo/go-kill/common/app"
	"github.com/linggaaskaedo/go-kill/common/database"
	"github.com/linggaaskaedo/go-kill/common/logger"
	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/linggaaskaedo/go-kill/common/server"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/config"
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
	log, err := logger.New(cfg.Logger)
	if err != nil {
		panic(err)
	}

	log.Info().Msg("Starting user service...")

	// Create application with options
	application := app.New(app.WithShutdownTimeout(15*time.Second), app.WithLogger(log))

	// Connect to database
	db0, err := database.New(cfg.Database["database-0"])
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	application.Add(app.NewDatabaseComponent(db0, log))

	// Query Loader Initialization
	queryLoader := query.New(cfg.Query)
	application.Add(app.NewQueryComponent(queryLoader, log))

	// HTTP Server
	httpServer := server.New(cfg.Server)
	// router := httpServer.Engine()
	// v1 := router.Group("/api/v1")
	// {
	// 	v1.GET("/health", h.HealthCheck)
	// 	users := v1.Group("/users")
	// 	{
	// 		users.POST("", h.CreateUser)
	// 	}
	// }
	application.Add(app.NewHTTPServerComponent(httpServer, log))

	// Run
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}
}
