// internal/services/transaction_service_test.go
package services_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"avito-shop/internal/models"
	"avito-shop/internal/repositories"
	"avito-shop/internal/services"
)

func (m *MockTransactionRepository) ToRepositoryType() *repositories.TransactionRepository {
	return (*repositories.TransactionRepository)(unsafe.Pointer(m))
}

func (m *MockUserRepository) ToRepositoryType() *repositories.UserRepository {
	return (*repositories.UserRepository)(unsafe.Pointer(m))
}

func TestTransactionService_TransferCoins(t *testing.T) {
	t.Run("Успешный перевод", func(t *testing.T) {
		mockDB, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

		mockUserRepo := &MockUserRepository{}
		mockTransactionRepo := &MockTransactionRepository{}

		sender := &models.User{
			Username: "sender",
			Coins:    500,
		}
		receiver := &models.User{
			Username: "receiver",
			Coins:    1000,
		}

		mockUserRepo.On("GetByUsername", "sender").Return(sender, nil)
		mockUserRepo.On("GetByUsername", "receiver").Return(receiver, nil)
		mockUserRepo.On("UpdateCoins", "sender", -100).Return(nil)
		mockUserRepo.On("UpdateCoins", "receiver", 100).Return(nil)
		mockTransactionRepo.On("Create", mock.Anything).Return(nil)

		transactionService := services.NewTransactionService(
			mockDB,
			mockTransactionRepo.ToRepositoryType(),
			mockUserRepo.ToRepositoryType(),
		)

		err := transactionService.TransferCoins("sender", "receiver", 100)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
}
