package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	UserID int `json:"user_id"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.UserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be a positive integer"})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "userID": user.ID, "coins": user.Coins})
}
