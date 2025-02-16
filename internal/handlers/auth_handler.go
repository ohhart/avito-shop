package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.AuthHandler = &AuthHandler{}

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors.ErrInvalidRequest.Error()})
		return
	}

	token, err := h.authService.AuthenticateOrRegister(req.Username, req.Password)
	if err != nil {
		switch {
		case err == errors.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Неверные учетные данные"})
		case err == errors.ErrInternalServer:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Внутренняя ошибка сервера"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Неизвестная ошибка"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
