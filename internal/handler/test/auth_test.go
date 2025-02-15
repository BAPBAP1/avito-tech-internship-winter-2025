package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct {
	LoginFunc func(userID int) (string, error)
}

func (m *MockAuthService) Login(userID int) (string, error) {
	return m.LoginFunc(userID)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockAuthService := &MockAuthService{
		LoginFunc: func(userID int) (string, error) {
			return "mock_token", nil
		},
	}
	authHandler := handler.NewAuthHandler(mockAuthService)
	router.POST("/auth", authHandler.Login)

	tests := []struct {
		name       string
		userID     int
		wantStatus int
	}{
		{"ValidUser", 1, http.StatusOK},
		{"InvalidUser", -1, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(map[string]int{"user_id": tt.userID})
			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "token")
			}
		})
	}
}
