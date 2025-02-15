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

type MockMerchService struct {
 ListMerchFunc     func() ([]model.Merch, error)
 PurchaseMerchFunc func(userID int, itemName string) error
}

func (m *MockMerchService) ListMerch() ([]model.Merch, error) {
 return m.ListMerchFunc()
}

func (m *MockMerchService) PurchaseMerch(userID int, itemName string) error {
 return m.PurchaseMerchFunc(userID, itemName)
}

func TestListMerch(t *testing.T) {
 gin.SetMode(gin.TestMode)
 router := gin.Default()

 mockMerchService := &MockMerchService{
  ListMerchFunc: func() ([]model.Merch, error) {
   return []model.Merch{{Name: "t-shirt", Price: 80}}, nil
  },
 }
 merchHandler := handler.NewMerchHandler(mockMerchService)
 router.GET("/merch", merchHandler.ListMerch)

 req := httptest.NewRequest(http.MethodGet, "/merch", nil)
 w := httptest.NewRecorder()

 router.ServeHTTP(w, req)

 assert.Equal(t, http.StatusOK, w.Code)
 var merchItems []model.Merch
 json.Unmarshal(w.Body.Bytes(), &merchItems)
 assert.NotEmpty(t, merchItems)
}

func TestPurchaseMerch(t *testing.T) {
 gin.SetMode(gin.TestMode)
 router := gin.Default()

 mockMerchService := &MockMerchService{
  PurchaseMerchFunc: func(userID int, itemName string) error {
   if itemName == "" {
    return service.ErrMerchNotFound
   }
   return nil
  },
 }
 merchHandler := handler.NewMerchHandler(mockMerchService)
 router.POST("/purchase", merchHandler.PurchaseMerch)

 tests := []struct {
  name       string
  body       map[string]string
  wantStatus int
 }{
  {"ValidPurchase", map[string]string{"item_name": "t-shirt"}, http.StatusOK},
  {"InvalidPurchaseEmptyItem", map[string]string{"item_name": ""}, http.StatusNotFound},
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   reqBody, _ := json.Marshal(tt.body)
   req := httptest.NewRequest(http.MethodPost, "/purchase", bytes.NewBuffer(reqBody))
   req.Header.Set("Authorization", "Bearer mocktoken")
   w := httptest.NewRecorder()

   router.ServeHTTP(w, req)

   assert.Equal(t, tt.wantStatus, w.Code)
  })
 }
}
