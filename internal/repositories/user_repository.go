package repositories

import (
	"gorm.io/gorm"

	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.UserRepository = &UserRepository{}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DB() *gorm.DB {
	return r.db
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, errors.ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateCoins(username string, amount int) error {
	return r.db.Model(&models.User{}).Where("username = ?", username).
		Update("coins", gorm.Expr("coins + ?", amount)).Error
}
