package services

import (
	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.ItemService = &ItemService{}

type ItemService struct {
	itemRepo ports.ItemRepository
}

func NewItemService(itemRepo ports.ItemRepository) *ItemService {
	return &ItemService{itemRepo: itemRepo}
}

func (s *ItemService) GetItemByName(name string) (*models.Item, error) {
	item, err := s.itemRepo.GetByName(name)
	if err != nil {
		return nil, errors.ErrItemNotFound
	}
	return item, nil
}

func (s *ItemService) GetAllItems() ([]models.Item, error) {
	return s.itemRepo.GetAll()
}
