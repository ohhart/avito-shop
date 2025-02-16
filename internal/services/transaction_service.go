package services

import (
	"log"
	"os"

	"gorm.io/gorm"

	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/repositories"
)

type TransactionService struct {
	db              *gorm.DB
	transactionRepo *repositories.TransactionRepository
	userRepo        *repositories.UserRepository
	logger          *log.Logger
}

func NewTransactionService(db *gorm.DB, transactionRepo *repositories.TransactionRepository, userRepo *repositories.UserRepository) *TransactionService {
	return &TransactionService{
		db:              db,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		logger:          log.New(os.Stdout, "TransactionService: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (s *TransactionService) TransferCoins(fromUsername string, toUsername string, amount int) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		fromUser, err := s.userRepo.GetByUsername(fromUsername)
		if err != nil {
			s.logger.Printf("Error getting sender (username: %s): %v", fromUsername, err)
			return errors.ErrInternalServer
		}

		toUser, err := s.userRepo.GetByUsername(toUsername)
		if err != nil {
			s.logger.Printf("Error getting receiver (username: %s): %v", toUsername, err)
			return errors.ErrInvalidRequest
		}

		if fromUser.ID == toUser.ID {
			s.logger.Printf("Attempted self-transfer: username=%s", fromUsername)
			return errors.ErrInvalidRequest
		}

		if amount <= 0 {
			s.logger.Printf("Invalid amount attempted: %d by username=%s", amount, fromUsername)
			return errors.ErrInvalidRequest
		}

		if fromUser.Coins < amount {
			s.logger.Printf("Insufficient funds: username=%s, current=%d, required=%d", fromUsername, fromUser.Coins, amount)
			return errors.ErrInsufficientFunds
		}

		if err := tx.Model(fromUser).Update("coins", gorm.Expr("coins - ?", amount)).Error; err != nil {
			s.logger.Printf("Failed to update sender balance (username=%s): %v", fromUsername, err)
			return errors.ErrInternalServer
		}

		if err := tx.Model(toUser).Update("coins", gorm.Expr("coins + ?", amount)).Error; err != nil {
			s.logger.Printf("Failed to update receiver balance (username=%s): %v", toUsername, err)
			return errors.ErrInternalServer
		}

		transaction := &models.Transaction{
			FromUserID: fromUser.ID,
			ToUserID:   toUser.ID,
			Amount:     amount,
		}

		if err := tx.Create(transaction).Error; err != nil {
			s.logger.Printf("Failed to create transaction record: %v", err)
			return errors.ErrInternalServer
		}

		return nil
	})

	if err != nil {
		s.logger.Printf("Transaction failed: %v", err)
		return err
	}

	return nil
}
