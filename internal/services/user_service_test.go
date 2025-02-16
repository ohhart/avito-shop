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

func TestUserService_GetUserInfo(t *testing.T) {
	t.Run("Успешное получение информации", func(t *testing.T) {
		mockDB, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
		mockUserRepo := new(MockUserRepository)
		mockInventoryRepo := new(MockInventoryRepository)
		mockTransactionRepo := new(MockTransactionRepository)

		user := &models.User{
			Username: "testuser",
			Coins:    500,
		}

		mockUserRepo.On("GetByUsername", "testuser").Return(user, nil)
		mockInventoryRepo.On("GetUserInventory", mock.Anything).Return([]models.Inventory{}, nil)
		mockTransactionRepo.On("GetUserTransactionHistory", mock.Anything).Return([]models.Transaction{}, nil)

		userService := services.NewUserService(
			mockDB,
			mockUserRepo,
			mockInventoryRepo,
			mockTransactionRepo,
		)

		info, err := userService.GetUserInfo("testuser")

		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, 500, info.Coins)
	})
}
