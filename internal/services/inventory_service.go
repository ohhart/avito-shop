package services

import (
	"gorm.io/gorm"

	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.InventoryService = &InventoryService{}

type InventoryService struct {
	db            *gorm.DB
	inventoryRepo ports.InventoryRepository
	userRepo      ports.UserRepository
	itemRepo      ports.ItemRepository
}

func NewInventoryService(
	db *gorm.DB,
	inventoryRepo ports.InventoryRepository,
	userRepo ports.UserRepository,
	itemRepo ports.ItemRepository,
) *InventoryService {
	return &InventoryService{
		db:            db,
		inventoryRepo: inventoryRepo,
		userRepo:      userRepo,
		itemRepo:      itemRepo,
	}
}

func (s *InventoryService) BuyItem(username, itemName string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepo.GetByUsername(username)
		if err != nil {
			return errors.ErrNotFound
		}

		item, err := s.itemRepo.GetByName(itemName)
		if err != nil {
			return errors.ErrItemNotFound
		}

		if user.Coins < item.Price {
			return errors.ErrInsufficientFunds
		}

		if err := s.userRepo.UpdateCoins(username, -item.Price); err != nil {
			return errors.ErrInternalServer
		}

		if err := s.inventoryRepo.AddItem(user.ID, itemName); err != nil {
			return errors.ErrInternalServer
		}

		return nil
	})
}

func (s *InventoryService) GetUserInventory(username string) ([]models.Inventory, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	return s.inventoryRepo.GetUserInventory(user.ID)
}
