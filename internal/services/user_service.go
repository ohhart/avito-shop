package services

import (
	"gorm.io/gorm"

	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.UserService = &UserService{}

type UserService struct {
	db              *gorm.DB
	userRepo        ports.UserRepository
	inventoryRepo   ports.InventoryRepository
	transactionRepo ports.TransactionRepository
}

func NewUserService(
	db *gorm.DB,
	userRepo ports.UserRepository,
	inventoryRepo ports.InventoryRepository,
	transactionRepo ports.TransactionRepository,
) *UserService {
	return &UserService{
		db:              db,
		userRepo:        userRepo,
		inventoryRepo:   inventoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *UserService) GetUserInfo(username string) (*models.InfoResponse, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventoryRepo.GetUserInventory(user.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	transactions, err := s.transactionRepo.GetUserTransactionHistory(user.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return s.prepareUserInfo(user, inventory, transactions)
}

func (s *UserService) TransferCoins(fromUsername, toUsername string, amount int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		fromUser, err := s.userRepo.GetByUsername(fromUsername)
		if err != nil {
			return errors.ErrNotFound
		}

		toUser, err := s.userRepo.GetByUsername(toUsername)
		if err != nil {
			return errors.ErrNotFound
		}

		if fromUser.ID == toUser.ID {
			return errors.ErrInvalidRequest
		}

		if amount <= 0 {
			return errors.ErrInvalidRequest
		}

		if fromUser.Coins < amount {
			return errors.ErrInsufficientFunds
		}

		if err := s.userRepo.UpdateCoins(fromUsername, -amount); err != nil {
			return errors.ErrInternalServer
		}

		if err := s.userRepo.UpdateCoins(toUsername, amount); err != nil {
			return errors.ErrInternalServer
		}

		transaction := &models.Transaction{
			FromUserID: fromUser.ID,
			ToUserID:   toUser.ID,
			Amount:     amount,
		}

		if err := s.transactionRepo.Create(transaction); err != nil {
			return errors.ErrInternalServer
		}

		return nil
	})
}

func (s *UserService) prepareUserInfo(
	user *models.User,
	inventory []models.Inventory,
	transactions []models.Transaction,
) (*models.InfoResponse, error) {
	inventoryResponse := make([]models.InventoryResponse, len(inventory))
	for i, item := range inventory {
		inventoryResponse[i] = models.InventoryResponse{
			Type:     item.ItemName,
			Quantity: item.Quantity,
		}
	}

	coinHistory, err := s.prepareCoinHistory(user.ID, transactions)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &models.InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventoryResponse,
		CoinHistory: coinHistory,
	}, nil
}

func (s *UserService) prepareCoinHistory(userID uint, transactions []models.Transaction) (models.CoinHistory, error) {
	coinHistory := models.CoinHistory{
		Received: []models.TransactionHistory{},
		Sent:     []models.TransactionHistory{},
	}

	for _, t := range transactions {
		if t.FromUserID == userID {
			toUser, err := s.userRepo.GetByID(t.ToUserID)
			if err != nil {
				return models.CoinHistory{}, errors.ErrInternalServer
			}
			coinHistory.Sent = append(coinHistory.Sent, models.TransactionHistory{
				ToUser: toUser.Username,
				Amount: t.Amount,
			})
		} else {
			fromUser, err := s.userRepo.GetByID(t.FromUserID)
			if err != nil {
				return models.CoinHistory{}, errors.ErrInternalServer
			}
			coinHistory.Received = append(coinHistory.Received, models.TransactionHistory{
				FromUser: fromUser.Username,
				Amount:   t.Amount,
			})
		}
	}

	return coinHistory, nil
}
