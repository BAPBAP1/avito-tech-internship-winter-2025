package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/service"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

type TransferRequest struct {
	ReceiverID int `json:"receiver_id"`
	Amount     int `json:"amount"`
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	senderID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req TransferRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.ReceiverID <= 0 || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver ID and Amount must be positive"})
		return
	}

	if int(senderID.(float64)) == req.ReceiverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer coins to yourself"})
		return
	}

	err := h.walletService.Transfer(c.Request.Context(), int(senderID.(float64)), req.ReceiverID, req.Amount)
	if err == service.ErrInsufficientFunds {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	} else if err == service.ErrInvalidAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer amount"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer coins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Coins transferred successfully"})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	wallet, err := h.walletService.GetWallet(c.Request.Context(), int(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallet"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) GetWalletHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	history, err := h.walletService.GetWalletHistory(c.Request.Context(), int(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallet history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
