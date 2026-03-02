package rest

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (e *rest) handleLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_request_body")
		e.httpRespError(c, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	loginReq := &authpb.LoginRequest{
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

	refreshReq := &authpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	resp, err := e.svc.Auth.RefreshToken(ctx, refreshReq)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleLogout(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPUnauthorized, "no_authorization_header"))
		return
	}

	// Remove "Bearer " prefix
	token := authHeader[7:]

	// Parse token to get user_id
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if err != nil || !parsedToken.Valid {
		e.httpRespError(c, x.NewWithCode(x.CodeHTTPUnauthorized, "invalid_token"))
		return
	}

	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Call internal Logout method
	logoutReq := &authpb.LogoutRequest{
		Token:  token,
		UserId: userID,
	}

	resp, _ := e.svc.Auth.Logout(ctx, logoutReq)

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}

func (e *rest) handleHealth(c *gin.Context) {
	resp := map[string]string{
		"status":  "healthy",
		"service": "auth-service",
	}

	e.httpRespSuccess(c, http.StatusOK, resp, nil)
}
