package handler

import (
 "net/http"

 "github.com/gin-gonic/gin"

 "github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/repository"
)

type UserHandler struct {
 userRepo        *repository.UserRepository
 transactionRepo *repository.TransactionRepository 
}

func NewUserHandler(userRepo *repository.UserRepository, transactionRepo *repository.TransactionRepository) *UserHandler {
 return &UserHandler{userRepo: userRepo, transactionRepo: transactionRepo}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
 userID, exists := c.Get("userID")
 if !exists {
  c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
  return
 }

 user, err := h.userRepo.GetByID(c.Request.Context(), int(userID.(float64)))
 if err != nil {
  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
  return
 }

 c.JSON(http.StatusOK, gin.H{
  "id":    user.ID,
  "coins": user.Coins,
 })
}
