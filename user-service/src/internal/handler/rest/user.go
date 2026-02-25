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
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	user, err := e.svc.User.RegisterUser(c.Request.Context(), req, e.grpc)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPInternalServerError, "invalid_request_body"))
		return
	}

	// c.JSON(http.StatusOK, user)
	e.httpRespSuccess(c, http.StatusCreated, user, nil)
}
