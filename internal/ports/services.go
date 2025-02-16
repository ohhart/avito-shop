package ports

import "avito-shop/internal/models"

type AuthService interface {
	AuthenticateOrRegister(username, password string) (string, error)
}

type UserService interface {
	GetUserInfo(username string) (*models.InfoResponse, error)
	TransferCoins(fromUsername, toUsername string, amount int) error
}

type InventoryService interface {
	BuyItem(username, itemName string) error
	GetUserInventory(username string) ([]models.Inventory, error)
}

type ItemService interface {
	GetItemByName(name string) (*models.Item, error)
	GetAllItems() ([]models.Item, error)
}

type TransactionService interface {
	TransferCoins(fromUsername, toUsername string, amount int) error
}
