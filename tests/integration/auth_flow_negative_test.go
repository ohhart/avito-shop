package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"avito-shop/internal/models"
)

// Неверное имя пользователя
func TestAuthHandlerInvalidUsername(t *testing.T) {
	application, _ := setupTestApp()

	authReq := models.AuthRequest{
		Username: "nonexistentuser",
		Password: "password123",
	}
	authBody, _ := json.Marshal(authReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	// Изменяем ожидание, так как сервис регистрирует нового пользователя
	assert.Equal(t, http.StatusOK, w.Code)
}

// Неверный пароль
func TestAuthHandlerInvalidPassword(t *testing.T) {
	application, _ := setupTestApp()
	//Сначала регистрируем
	authReq := models.AuthRequest{
		Username: "testuserAuth",
		Password: "testpassAuth",
	}
	authBody, _ := json.Marshal(authReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	authReq = models.AuthRequest{
		Username: "testuserAuth",
		Password: "wrongpassword",
	}
	authBody, _ = json.Marshal(authReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // Ожидаем 401
}

// Пустое тело запроса
func TestAuthHandlerEmptyBody(t *testing.T) {
	application, _ := setupTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", nil) // Пустое тело
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400
}

// Пустое имя пользователя
func TestAuthHandlerEmptyUsername(t *testing.T) {
	application, _ := setupTestApp()

	authReq := models.AuthRequest{
		Username: "",
		Password: "password123",
	}
	authBody, _ := json.Marshal(authReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Пустой пароль
func TestAuthHandlerEmptyPassword(t *testing.T) {
	application, _ := setupTestApp()

	authReq := models.AuthRequest{
		Username: "testuser",
		Password: "",
	}
	authBody, _ := json.Marshal(authReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
