package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"avito-shop/internal/models"
	"avito-shop/internal/services"
)

func TestInventoryService_BuyItem(t *testing.T) {
	t.Run("Успешная покупка", func(t *testing.T) {
		mockDB, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
		mockUserRepo := new(MockUserRepository)
		mockItemRepo := new(MockItemRepository)
		mockInventoryRepo := new(MockInventoryRepository)

		user := &models.User{
			Username: "testuser",
			Coins:    1000,
		}

		item := &models.Item{
			Name:  "t-shirt",
			Price: 80,
		}

		mockUserRepo.On("GetByUsername", "testuser").Return(user, nil)
		mockItemRepo.On("GetByName", "t-shirt").Return(item, nil)
		mockUserRepo.On("UpdateCoins", "testuser", -80).Return(nil)
		mockInventoryRepo.On("AddItem", mock.Anything, "t-shirt").Return(nil)

		inventoryService := services.NewInventoryService(
			mockDB,
			mockInventoryRepo,
			mockUserRepo,
			mockItemRepo,
		)

		err := inventoryService.BuyItem("testuser", "t-shirt")

		assert.NoError(t, err)
	})
}
