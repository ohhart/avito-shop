package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-shop/internal/errors"
	"avito-shop/internal/services"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) TransferCoins(c *gin.Context) {
	var req struct {
		ToUsername string `json:"toUser" binding:"required"`
		Amount     int    `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.transactionService.TransferCoins(username.(string), req.ToUsername, req.Amount)
	if err != nil {
		switch err {
		case errors.ErrInvalidRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		case errors.ErrInsufficientFunds:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Перевод успешен"})
}
