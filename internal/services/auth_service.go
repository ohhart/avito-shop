package services

import (
	"avito-shop/internal/auth"
	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/ports"
)

var _ ports.AuthService = &AuthService{}

type AuthService struct {
	userRepo ports.UserRepository
}

func NewAuthService(userRepo ports.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) AuthenticateOrRegister(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		// Если пользователя нет, создаем нового
		if err == errors.ErrUserNotFound {
			hashedPassword, hashErr := auth.HashPassword(password)
			if hashErr != nil {
				return "", errors.ErrInternalServer
			}

			newUser := &models.User{
				Username:     username,
				PasswordHash: hashedPassword,
				Coins:        1000,
			}

			err = s.userRepo.Create(newUser)
			if err != nil {
				return "", err
			}

			return auth.GenerateJWT(newUser.Username)
		}
		return "", err
	}

	// Проверяем пароль
	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.ErrInvalidCredentials
	}

	// Генерируем токен
	return auth.GenerateJWT(user.Username)
}
