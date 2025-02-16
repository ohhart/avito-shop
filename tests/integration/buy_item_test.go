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

func TestBuyItemIntegration(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}

	// Регистрация пользователя
	authReq := models.AuthRequest{
		Username: "testuser",
		Password: "testpass",
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

	// Покупка товара
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/buy/t-shirt", nil)
	assert.NoError(t, err, "Ошибка при создании запроса /api/buy/t-shirt")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка инвентаря
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/info", nil)
	assert.NoError(t, err, "Ошибка при создании запроса /api/info")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	var infoResp models.InfoResponse
	err = json.Unmarshal(w.Body.Bytes(), &infoResp)
	assert.NoError(t, err, "Ошибка при unmarshal infoResp")

	assert.Less(t, infoResp.Coins, 1000)
	assert.Contains(t, infoResp.Inventory, models.InventoryResponse{
		Type:     "t-shirt",
		Quantity: 1,
	})
}
