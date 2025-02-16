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

func TestUserInfoRetrieval(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}

	// Регистрация пользователя
	authReq := models.AuthRequest{
		Username: "infouser",
		Password: "password123",
	}
	authBody, err := json.Marshal(authReq)
	assert.NoError(t, err, "Ошибка при маршалинге authReq")

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	assert.NoError(t, err, "Ошибка при создании запроса /api/auth")
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var authResp struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err, "Ошибка при unmarshal authResp")

	// Получение информации о пользователе
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/info", nil)
	assert.NoError(t, err, "Ошибка при создании запроса /api/info")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var userInfo models.InfoResponse
	err = json.Unmarshal(w.Body.Bytes(), &userInfo)
	assert.NoError(t, err, "Ошибка при unmarshal userInfo")

	assert.Equal(t, 1000, userInfo.Coins)
	assert.Empty(t, userInfo.Inventory)
	assert.Empty(t, userInfo.CoinHistory.Received)
	assert.Empty(t, userInfo.CoinHistory.Sent)
}
