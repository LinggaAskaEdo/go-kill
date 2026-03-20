package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openpcc/openpcc/uuidv7"
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

	if productID == "" {
		e.httpRespError(c, errProductIDRequired)
		return
	}
	if _, err := uuidv7.Parse(productID); err != nil {
		e.httpRespError(c, errInvalidProductIDFormat)
		return
	}

	resp, err := e.svc.Product.GetProduct(ctx, productID)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	if resp == nil {
		e.httpRespError(c, errProductNotFound)
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

	if productID == "" {
		e.httpRespError(c, errProductIDRequired)
		return
	}
	if _, err := uuidv7.Parse(productID); err != nil {
		e.httpRespError(c, errInvalidProductIDFormat)
		return
	}

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

	if categoryID == "" {
		e.httpRespError(c, errCategoryIDRequired)
		return
	}
	if _, err := uuidv7.Parse(categoryID); err != nil {
		e.httpRespError(c, errInvalidCategoryIDFormat)
		return
	}

	resp, err := e.svc.Product.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}
