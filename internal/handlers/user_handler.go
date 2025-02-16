package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-shop/internal/errors"
	"avito-shop/internal/ports"
)

var _ ports.UserHandler = &UserHandler{}

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Неавторизованный доступ"})
		return
	}

	info, err := h.userService.GetUserInfo(username.(string))
	if err != nil {
		switch {
		case err == errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"errors": "Пользователь не найден"})
		case err == errors.ErrInternalServer:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Внутренняя ошибка сервера"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Неизвестная ошибка"})
		}
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *UserHandler) SendCoins(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Неавторизованный доступ"})
		return
	}

	var req struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Неверный запрос"})
		return
	}

	err := h.userService.TransferCoins(username.(string), req.ToUser, req.Amount)
	if err != nil {
		switch {
		case err == errors.ErrInsufficientFunds:
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Недостаточно средств"})
		case err == errors.ErrNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Пользователь не найден"})
		case err == errors.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Некорректный запрос"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Внутренняя ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Монеты успешно переведены"})
}
