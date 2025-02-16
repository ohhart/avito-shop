package repositories

import (
	"gorm.io/gorm"

	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.TransactionRepository = &TransactionRepository{}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetUserTransactionHistory(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}
