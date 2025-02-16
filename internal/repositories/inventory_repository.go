package repositories

import (
	"gorm.io/gorm"

	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.InventoryRepository = &InventoryRepository{}

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetUserInventory(userID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	if err := r.db.Where("user_id = ?", userID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *InventoryRepository) AddItem(userID uint, itemName string) error {
	var inventory models.Inventory
	err := r.db.Where("user_id = ? AND item_name = ?", userID, itemName).First(&inventory).Error
	if err == gorm.ErrRecordNotFound {
		inventory = models.Inventory{UserID: userID, ItemName: itemName, Quantity: 1}
		return r.db.Create(&inventory).Error
	} else if err != nil {
		return err
	}

	return r.db.Model(&inventory).Update("quantity", gorm.Expr("quantity + ?", 1)).Error
}
