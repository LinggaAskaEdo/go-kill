package rest

import (
	"context"
	"net/http"
	"sync"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"

	"github.com/gin-gonic/gin"
)

var onceRestHandler = &sync.Once{}

type rest struct {
	gin *gin.Engine
	svc *service.Service
}

func InitRestHandler(gin *gin.Engine, svc *service.Service) {
	var e *rest

	onceRestHandler.Do(func() {
		e = &rest{
			gin: gin,
			svc: svc,
		}

		e.Serve()
	})
}

func (e *rest) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			c.Abort()
			return
		}

		// Validate with Auth Service
		resp, err := e.svc.User.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{Token: authHeader[7:]})
		if err != nil || !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_auth_id", resp.UserId)
		c.Set("email", resp.Email)
		c.Next()
	}
}

func (e *rest) Serve() {
	// User
	e.gin.POST("/api/v1/users/register", e.handleRegister)
	e.gin.GET("/api/v1/users/me", e.authMiddleware(), e.handleGetMe)
	e.gin.GET("/api/v1/users/me/activities", e.authMiddleware(), e.handleGetActivities)
	e.gin.GET("/api/v1/users/me/addresses", e.authMiddleware(), e.handleGetAddresses)
	e.gin.POST("/api/v1/users/me/addresses", e.authMiddleware(), e.handleCreateAddress)
}
