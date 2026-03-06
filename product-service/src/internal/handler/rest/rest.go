package rest

import (
	"sync"

	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service"

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

func (e *rest) Serve() {
	e.gin.GET("/api/v1/products", e.handleListProducts)
	e.gin.GET("/api/v1/products/:id", e.handleGetProduct)
	e.gin.GET("/api/v1/categories", e.handleListCategories)
	e.gin.GET("/api/v1/products/:id/categories", e.handleGetCategoriesByProduct)
	e.gin.GET("/api/v1/categories/:id/products", e.handleGetProductsByCategory)
}
