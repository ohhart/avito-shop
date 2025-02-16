package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-shop/internal/errors"
	"avito-shop/internal/ports"
)

var _ ports.InventoryHandler = &InventoryHandler{}

type InventoryHandler struct {
	inventoryService ports.InventoryService
}

func NewInventoryHandler(inventoryService ports.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) BuyItem(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Неавторизованный доступ"})
		return
	}

	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Не указан товар"})
		return
	}

	err := h.inventoryService.BuyItem(username.(string), itemName)
	if err != nil {
		switch {
		case err == errors.ErrInsufficientFunds:
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Недостаточно средств"})
		case err == errors.ErrItemNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Товар не найден"})
		case err == errors.ErrInternalServer:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Внутренняя ошибка сервера"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Неизвестная ошибка"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно куплен"})
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Неавторизованный доступ"})
		return
	}

	inventory, err := h.inventoryService.GetUserInventory(username.(string))
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

	c.JSON(http.StatusOK, inventory)
}
