package ports

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	Login(c *gin.Context)
}

type UserHandler interface {
	GetUserInfo(c *gin.Context)
	SendCoins(c *gin.Context)
}

type InventoryHandler interface {
	BuyItem(c *gin.Context)
	GetInventory(c *gin.Context)
}
