package rest

import (
	"sync"

	rpc "github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"

	"github.com/gin-gonic/gin"
)

var onceRestHandler = &sync.Once{}

type rest struct {
	gin  *gin.Engine
	svc  *service.Service
	grpc *rpc.Grpc
}

func InitRestHandler(gin *gin.Engine, svc *service.Service, grpc *rpc.Grpc) {
	var e *rest

	onceRestHandler.Do(func() {
		e = &rest{
			gin:  gin,
			svc:  svc,
			grpc: grpc,
		}

		e.Serve()
	})
}
func (e *rest) Serve() {
	// User
	e.gin.POST("/api/v1/users/register", e.handleRegister)
	// e.gin.GET("/api/v1/users/me", e.authMiddleware(), e.handleGetMe)
	// e.gin.GET("/api/v1/users/me/activities", e.authMiddleware(), e.handleGetActivities)
	// e.gin.GET("/api/v1/users/me/addresses", e.authMiddleware(), e.handleGetAddresses)
	// e.gin.POST("/api/v1/users/me/addresses", e.authMiddleware(), e.handleCreateAddress)
}
