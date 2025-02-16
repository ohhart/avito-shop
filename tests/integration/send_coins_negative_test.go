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
func TestSendCoinsInvalidToken(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "user2",
		Amount: 50,
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token") // Невалидный токен
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Несуществующий получатель
func TestSendCoinsNonexistentRecipient(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend1",
		Password: "testpassSend1",
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

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "nonexistentuser",
		Amount: 50,
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400
}

// Отрицательная сумма
func TestSendCoinsNegativeAmount(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend2",
		Password: "testpassSend2",
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
	//Регистрируем 2
	authReq2 := models.AuthRequest{
		Username: "testuserSend3",
		Password: "testpassSend3",
	}
	authBody2, _ := json.Marshal(authReq2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody2))
	req2.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w2, req2)

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "testuserSend3",
		Amount: -50, // Отрицательная сумма
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400
}

// Недостаточно средств
func TestSendCoinsInsufficientFunds(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend4",
		Password: "testpassSend4",
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
	//Регистрируем 2
	authReq2 := models.AuthRequest{
		Username: "testuserSend5",
		Password: "testpassSend5",
	}
	authBody2, _ := json.Marshal(authReq2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody2))
	req2.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w2, req2)

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "testuserSend5",
		Amount: 2000, // Больше, чем есть на балансе
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400
}

// Перевод самому себе
func TestSendCoinsToSelf(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend6",
		Password: "testpassSend6",
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

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "testuserSend6", // Самому себе
		Amount: 50,
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) //Ожидаем ошибку
}

// Пустое поле ToUser
func TestSendCoinsEmptyToUser(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend7",
		Password: "testpassSend7",
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

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "", // Пустое поле
		Amount: 50,
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Пустое поле Amount
func TestSendCoinsEmptyAmount(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend8",
		Password: "testpassSend8",
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
	//Регистрируем 2
	authReq2 := models.AuthRequest{
		Username: "testuserSend9",
		Password: "testpassSend9",
	}
	authBody2, _ := json.Marshal(authReq2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody2))
	req2.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w2, req2)
	sendCoinReq := map[string]interface{}{ // Используем map, чтобы не создавать структуру
		"toUser": "testuserSend9",
		"amount": nil, // Пустое поле
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Неверный Content-Type
func TestSendCoinsInvalidContentType(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend10",
		Password: "testpassSend10",
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
	//Регистрируем 2
	authReq2 := models.AuthRequest{
		Username: "testuserSend11",
		Password: "testpassSend11",
	}
	authBody2, _ := json.Marshal(authReq2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody2))
	req2.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w2, req2)

	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "testuserSend11",
		Amount: 50,
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "text/plain") // Неверный Content-Type
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем ошибку
}

// Нулевая сумма перевода
func TestSendCoinsZeroAmount(t *testing.T) {
	application, _ := setupTestApp()
	//Регистрируем
	authReq := models.AuthRequest{
		Username: "testuserSend12",
		Password: "testpassSend12",
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
	//Регистрируем 2
	authReq2 := models.AuthRequest{
		Username: "testuserSend13",
		Password: "testpassSend13",
	}
	authBody2, _ := json.Marshal(authReq2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(authBody2))
	req2.Header.Set("Content-Type", "application/json")
	application.Router().ServeHTTP(w2, req2)
	sendCoinReq := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: "testuserSend13",
		Amount: 0, // Ноль
	}
	sendBody, _ := json.Marshal(sendCoinReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Ожидаем 400 (сумма должна быть > 0)
}
