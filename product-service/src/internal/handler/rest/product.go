package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *rest) handleListProducts(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := e.svc.Product.ListProduct(ctx)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleGetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("id")

	resp, err := e.svc.Product.GetProduct(ctx, productID)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleListCategories(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := e.svc.Product.ListCategories(ctx)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleGetCategoriesByProduct(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("id")

	resp, err := e.svc.Product.GetCategoriesByProduct(ctx, productID)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleGetProductsByCategory(c *gin.Context) {
	ctx := c.Request.Context()
	categoryID := c.Param("id")

	resp, err := e.svc.Product.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}
