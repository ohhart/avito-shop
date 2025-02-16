package ports

import (
	"gorm.io/gorm"

	"avito-shop/internal/models"
)

type UserRepository interface {
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	UpdateCoins(username string, amount int) error
	DB() *gorm.DB
}

type TransactionRepository interface {
	GetUserTransactionHistory(userID uint) ([]models.Transaction, error)
	Create(transaction *models.Transaction) error
}

type ItemRepository interface {
	GetByName(name string) (*models.Item, error)
	GetAll() ([]models.Item, error)
}

type InventoryRepository interface {
	GetUserInventory(userID uint) ([]models.Inventory, error)
	AddItem(userID uint, itemName string) error
}
