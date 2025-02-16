package repositories

import (
	"gorm.io/gorm"

	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.ItemRepository = &ItemRepository{}

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetByName(name string) (*models.Item, error) {
	var item models.Item
	if err := r.db.Where("name = ?", name).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) GetAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	return items, err
}
