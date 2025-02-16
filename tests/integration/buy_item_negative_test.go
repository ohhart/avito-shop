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

// Невалидный токен
func TestBuyItemInvalidToken(t *testing.T) {
	application, _ := setupTestApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/buy/t-shirt", nil)
	req.Header.Set("Authorization", "Bearer invalid_token") // Невалидный токен
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Несуществующий товар
func TestBuyItemNotFound(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем пользователя
	authReq := models.AuthRequest{
		Username: "testuserBuy",
		Password: "testpassBuy",
	}
	authBody, _ := json.Marshal(authReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var authResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &authResp)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/nonexistent-item", nil)
	req.Header.Set("Authorization", "Bearer "+authResp.Token) // Валидный токен
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400, т.к. товар не найден
}

// Недостаточно средств
func TestBuyItemInsufficientFunds(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем пользователя
	authReq := models.AuthRequest{
		Username: "testuserBuy2",
		Password: "testpassBuy2",
	}
	authBody, _ := json.Marshal(authReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var authResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &authResp)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/pink-hoody", nil) // Дорогой товар (500 монет)
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400
}

// Пустой item в URL
func TestBuyItemEmptyItemName(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserBuy3",
		Password: "testpassBuy3",
	}
	authBody, _ := json.Marshal(authReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var authResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &authResp)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/", nil) // Пустой item
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code) // 404, потому что нет такого маршрута
}

// Покупка после исчерпания средств (проверка состояния БД)
func TestBuyItemThenInsufficientFunds(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserBuy4",
		Password: "testpassBuy4",
	}
	authBody, _ := json.Marshal(authReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody))
	req.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w, req)

	var authResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &authResp)
	// Сначала покупаем дешевый товар
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/pen", nil) // 10 монет
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// ... и так 99 раз, чтобы потратить 990 монет ...
	for i := 0; i < 99; i++ {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/buy/pen", nil) // 10 монет
		req.Header.Set("Authorization", "Bearer "+authResp.Token)
		application.Router().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Теперь пытаемся купить товар за 20 монет
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/buy/cup", nil) // 20 монет
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400 (недостаточно средств)
}
