package handler

import (
 "net/http"
 "strconv"

 "github.com/gin-gonic/gin"

 "github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
 "github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/service"
)

type MerchHandler struct {
 merchService *service.MerchService
}

func NewMerchHandler(merchService *service.MerchService) *MerchHandler {
 return &MerchHandler{merchService: merchService}
}

func (h *MerchHandler) ListMerch(c *gin.Context) {
 merchItems, err := h.merchService.ListMerch(c.Request.Context())
 if err != nil {
  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list merch"})
  return
 }
 c.JSON(http.StatusOK, merchItems)
}

type PurchaseRequest struct {
 ItemName string `json:"item_name"`
}

func (h *MerchHandler) PurchaseMerch(c *gin.Context) {
 userID, exists := c.Get("userID")
 if !exists {
  c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
  return
 }

 var req PurchaseRequest
 if err := c.BindJSON(&req); err != nil {
  c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
  return
 }

 if req.ItemName == "" {
  c.JSON(http.StatusBadRequest, gin.H{"error": "Item name is required"})
  return
 }

 err := h.merchService.PurchaseMerch(c.Request.Context(), int(userID.(float64)), req.ItemName)
 if err == service.ErrInsufficientFunds {
  c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
  return
 } else if err == service.ErrMerchNotFound {
  c.JSON(http.StatusNotFound, gin.H{"error": "Merch item not found"})
  return
 } else if err != nil {
  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase merch"})
  return
 }

 c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Merch purchased successfully"})
}

func (h *MerchHandler) ListPurchases(c *gin.Context) {
 userID, exists := c.Get("userID")
 if !exists {
  c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
  return
 }

 purchases, err := h.merchService.ListPurchases(c.Request.Context(), int(userID.(float64)))
 if err != nil {
  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list purchases"})
  return
 }
 c.JSON(http.StatusOK, purchases)
}


func (h *MerchHandler) ListPurchasesByUserID(c *gin.Context) {
 userIDStr := c.Param("user_id")
 if userIDStr == "" {
  c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
  return
 }

 userID, err := strconv.Atoi(userIDStr)
 if err != nil {
  c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
  return
 }

 purchases, err := h.merchService.ListPurchases(c.Request.Context(), userID)
 if err != nil {
  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list purchases"})
  return
 }
 c.JSON(http.StatusOK, purchases)
}


func (h *MerchHandler) CreatePurchaseForUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	itemName := c.Param("item_name")
   
	if userIDStr == "" || itemName == "" {
	 c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Item Name are required"})
	 return
	}
   
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
	 c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
	 return
	}
   
	err = h.merchService.PurchaseMerch(c.Request.Context(), userID, itemName)
	if err != nil {
	 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase"})
	 return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Purchase created successfully"})
   }
   
   type PurchaseHistoryResponse struct {
	Purchases []model.Purchase `json:"purchases"`
   }
   
   