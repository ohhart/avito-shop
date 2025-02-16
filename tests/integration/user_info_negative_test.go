package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Невалидный токен
func TestGetUserInfoInvalidToken(t *testing.T) {
	application, err := setupTestApp()
	if err != nil {
		t.Fatalf("Ошибка при настройке тестового приложения: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	application.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
