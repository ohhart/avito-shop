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

func TestAuthenticationFlow(t *testing.T) {
	application, err := setupTestApp()
	assert.NoError(t, err)

	// Тест 1: Регистрация нового пользователя
	authReq := models.AuthRequest{
		Username: "newuser",
		Password: "password123",
	}
	authBody, err := json.Marshal(authReq)
	assert.NoError(t, err)

	// Первый запрос - регистрация
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	// Проверяем успешную регистрацию
	assert.Equal(t, http.StatusOK, w.Code)

	var firstResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &firstResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, firstResponse.Token)

	// Тест 2: Вход существующего пользователя
	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	// Проверяем успешный вход
	assert.Equal(t, http.StatusOK, w.Code)

	var secondResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &secondResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, secondResponse.Token)

	// Тест 3: Проверка доступа к защищенному эндпоинту
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/info", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+secondResponse.Token)
	application.Router().ServeHTTP(w, req)

	// Проверяем успешный доступ к защищенному ресурсу
	assert.Equal(t, http.StatusOK, w.Code)
}
