package rest

import (
	"net/http"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (e *rest) handleRegister(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_request_body")
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	resp, err := e.svc.User.RegisterUser(c.Request.Context(), req)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusCreated, resp, nil)
}

func (e *rest) handleGetMe(c *gin.Context) {
	ctx := c.Request.Context()
	userAuthID := c.GetString("user_auth_id")

	resp, err := e.svc.User.GetMe(ctx, userAuthID)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleGetActivities(c *gin.Context) {
	ctx := c.Request.Context()
	userAuthID := c.GetString("user_auth_id")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	resp, err := e.svc.User.GetActivities(ctx, userAuthID, page, limit)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleGetAddresses(c *gin.Context) {
	ctx := c.Request.Context()
	userAuthID := c.GetString("user_auth_id")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	resp, err := e.svc.User.GetAddresses(ctx, userAuthID, page, limit)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleCreateAddress(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateUserAddress
	userAuthID := c.GetString("user_auth_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_request_body")
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	resp, err := e.svc.User.CreateAddress(ctx, userAuthID, req)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}
