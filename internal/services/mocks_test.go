// internal/services/mocks_test.go
package services_test

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"avito-shop/internal/models"
)

// MockUserRepository имитирует интерфейс repositories.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateCoins(username string, amount int) error {
	args := m.Called(username, amount)
	return args.Error(0)
}

func (m *MockUserRepository) DB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// MockTransactionRepository имитирует интерфейс repositories.TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetUserTransactionHistory(userID uint) ([]models.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

// MockItemRepository имитирует интерфейс repositories.ItemRepository
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetByName(name string) (*models.Item, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

func (m *MockItemRepository) GetAll() ([]models.Item, error) {
	args := m.Called()
	return args.Get(0).([]models.Item), args.Error(1)
}

// MockInventoryRepository имитирует интерфейс repositories.InventoryRepository
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) GetUserInventory(userID uint) ([]models.Inventory, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) AddItem(userID uint, itemName string) error {
	args := m.Called(userID, itemName)
	return args.Error(0)
}
