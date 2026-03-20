package rest

import (
	"sync"

	rpc "github.com/linggaaskaedo/go-kill/auth-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"

	"github.com/gin-gonic/gin"
)

var onceRestHandler = &sync.Once{}

type rest struct {
	gin       *gin.Engine
	svc       *service.Service
	grpc      *rpc.Grpc
	jwtSecret []byte
}

func InitRestHandler(gin *gin.Engine, svc *service.Service, grpc *rpc.Grpc, jwtSecret string) {
	var e *rest

	onceRestHandler.Do(func() {
		e = &rest{
			gin:       gin,
			svc:       svc,
			grpc:      grpc,
			jwtSecret: []byte(jwtSecret),
		}

		e.Serve()
	})
}

func (e *rest) Serve() {
	e.gin.POST("/api/v1/auth/login", e.handleLogin)
	e.gin.POST("/api/v1/auth/refresh", e.handleRefresh)
	e.gin.POST("/api/v1/auth/logout", e.handleLogout)
	e.gin.GET("/health", e.handleHealth)
}
