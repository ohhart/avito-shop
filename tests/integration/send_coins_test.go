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

func TestSendCoinsIntegration(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}

	// Регистрация первого пользователя
	firstUser := models.AuthRequest{
		Username: "user1",
		Password: "pass1",
	}
	firstBody, err := json.Marshal(firstUser)
	assert.NoError(t, err, "Ошибка при маршалинге firstUser")

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(firstBody))
	assert.NoError(t, err, "Ошибка при создании запроса /api/auth для первого пользователя")
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var firstToken struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &firstToken)
	assert.NoError(t, err, "Ошибка при unmarshal firstToken")

	// Регистрация второго пользователя
	secondUser := models.AuthRequest{
		Username: "user2",
		Password: "pass2",
	}
	secondBody, err := json.Marshal(secondUser)
	assert.NoError(t, err, "Ошибка при маршалинге secondUser")

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(secondBody))
	assert.NoError(t, err, "Ошибка при создании запроса /api/auth для второго пользователя")
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var secondToken struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &secondToken)
	assert.NoError(t, err, "Ошибка при unmarshal secondToken")

	// Перевод монет
	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "user2",
		Amount: 50,
	}
	sendBody, err := json.Marshal(sendCoinReq)
	assert.NoError(t, err, "Ошибка при маршалинге sendCoinReq")

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	assert.NoError(t, err, "Ошибка при создании запроса /api/sendCoin")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+firstToken.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка баланса первого пользователя
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/info", nil)
	assert.NoError(t, err, "Ошибка при создании запроса /api/info для первого пользователя")
	req.Header.Set("Authorization", "Bearer "+firstToken.Token)
	application.Router().ServeHTTP(w, req)

	var firstUserInfo models.InfoResponse
	err = json.Unmarshal(w.Body.Bytes(), &firstUserInfo)
	assert.NoError(t, err, "Ошибка при unmarshal firstUserInfo")

	assert.Less(t, firstUserInfo.Coins, 1000)
	assert.Contains(t, firstUserInfo.CoinHistory.Sent, models.TransactionHistory{
		ToUser: "user2",
		Amount: 50,
	})

	// Проверка баланса второго пользователя
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/info", nil)
	assert.NoError(t, err, "Ошибка при создании запроса /api/info для второго пользователя")
	req.Header.Set("Authorization", "Bearer "+secondToken.Token)
	application.Router().ServeHTTP(w, req)

	var secondUserInfo models.InfoResponse
	err = json.Unmarshal(w.Body.Bytes(), &secondUserInfo)
	assert.NoError(t, err, "Ошибка при unmarshal secondUserInfo")

	assert.Contains(t, secondUserInfo.CoinHistory.Received, models.TransactionHistory{
		FromUser: "user1",
		Amount:   50,
	})
}
