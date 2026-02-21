package rest

import (
	"net/http"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"

	"github.com/gin-gonic/gin"
)

func (e *rest) handleRegister(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := e.svc.User.RegisterUser(c.Request.Context(), req, e.grpc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
