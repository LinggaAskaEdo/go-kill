package rest

import (
	"net/http"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

func (e *rest) handleLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_request_body")
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	loginReq := &dto.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		IpAddress: c.Request.Host,
		UserAgent: c.Request.UserAgent(),
	}

	resp, err := e.svc.Auth.Login(ctx, loginReq)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleRefresh(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_request_body")
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	resp, err := e.svc.Auth.RefreshToken(ctx, &req)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleLogout(c *gin.Context) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPUnauthorized, "no_authorization_header"))
		return
	}

	token := authHeader[7:]

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return e.jwtSecret, nil })
	if err != nil || !parsedToken.Valid {
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPUnauthorized, "invalid_token"))
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPUnauthorized, "invalid_token_claims"))
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPUnauthorized, "invalid_token_subject"))
		return
	}
	userID := sub

	logoutReq := &dto.LogoutRequest{
		Token:  token,
		UserId: userID,
	}

	resp, err := e.svc.Auth.Logout(ctx, logoutReq)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("logout_failed")
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPInternalServerError, "logout_failed"))
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleHealth(c *gin.Context) {
	resp := map[string]string{
		"status":  "healthy",
		"service": "auth-service",
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}
