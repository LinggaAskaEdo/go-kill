package http

import (
	"github.com/linggaaskaedo/go-kill/common/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Config struct {
	AppName string `yaml:"app_name"`
}

func Init(log zerolog.Logger, middleware middleware.Middleware, cfg Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Recovery())

	if cfg.AppName != "" {
		router.Use(otelgin.Middleware(cfg.AppName))
	}

	router.Use(middleware.Handler())
	router.Use(middleware.CORS())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	return router
}
