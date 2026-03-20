package rest

import (
	"sync"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	onceRestHandler = &sync.Once{}

	errProductIDRequired       = x.New("product ID is required")
	errInvalidProductIDFormat  = x.New("invalid product ID format")
	errProductNotFound         = x.NewWithCode(x.CodeSQLRecordDoesNotExist, "product not found")
	errCategoryIDRequired      = x.New("category ID is required")
	errInvalidCategoryIDFormat = x.New("invalid category ID format")
)

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
