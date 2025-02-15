package test

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/gin-gonic/gin"
 "github.com/stretchr/testify/assert"
 "github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/handler"
 "github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/service"
)

type MockWalletService struct {
 TransferFunc func(senderID, receiverID, amount int) error
}

func (m *MockWalletService) Transfer(senderID, receiverID, amount int) error {
 return m.TransferFunc(senderID, receiverID, amount)
}

func TestTransfer(t *testing.T) {
 gin.SetMode(gin.TestMode)
 router := gin.Default()

 mockWalletService := &MockWalletService{
  TransferFunc: func(senderID, receiverID, amount int) error {
   if amount <= 0 {
    return service.ErrInsufficientFunds
   }
   return nil
  },
 }
 walletHandler := handler.NewWalletHandler(mockWalletService)
 router.POST("/transfer", walletHandler.Transfer)

 tests := []struct {
  name       string
  body       map[string]int
  wantStatus int
 }{
  {"ValidTransfer", map[string]int{"receiver_id": 2, "amount": 100}, http.StatusOK},
  {"InvalidTransferNegativeAmount", map[string]int{"receiver_id": 2, "amount": -100}, http.StatusBadRequest},
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   reqBody, _ := json.Marshal(tt.body)
   req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(reqBody))
   req.Header.Set("Authorization", "Bearer mocktoken") 
   w := httptest.NewRecorder()

   router.ServeHTTP(w, req)

   assert.Equal(t, tt.wantStatus, w.Code)
  })
 }
}
