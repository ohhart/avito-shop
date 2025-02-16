// internal/services/auth_service_test.go
package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"avito-shop/internal/auth"
	"avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/services"
)

func TestAuthService_AuthenticateOrRegister(t *testing.T) {
	t.Run("Успешная регистрация нового пользователя", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByUsername", "newuser").Return(nil, errors.ErrUserNotFound)
		mockRepo.On("Create", mock.MatchedBy(func(user *models.User) bool {
			return user.Username == "newuser" &&
				user.Coins == 1000 &&
				user.PasswordHash != ""
		})).Return(nil)

		authService := services.NewAuthService(mockRepo)

		// Act
		token, err := authService.AuthenticateOrRegister("newuser", "password123")

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Успешная аутентификация существующего пользователя", func(t *testing.T) {
		// Arrange
		hashedPassword, _ := auth.HashPassword("password123")
		existingUser := &models.User{
			Username:     "existinguser",
			PasswordHash: hashedPassword,
		}

		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)

		authService := services.NewAuthService(mockRepo)

		// Act
		token, err := authService.AuthenticateOrRegister("existinguser", "password123")

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при неверном пароле", func(t *testing.T) {
		// Arrange
		hashedPassword, _ := auth.HashPassword("correctpassword")
		existingUser := &models.User{
			Username:     "existinguser",
			PasswordHash: hashedPassword,
		}

		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)

		authService := services.NewAuthService(mockRepo)

		// Act
		token, err := authService.AuthenticateOrRegister("existinguser", "wrongpassword")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidCredentials, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}
