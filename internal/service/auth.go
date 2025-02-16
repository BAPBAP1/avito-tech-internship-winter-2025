package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/repository"
)


type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, userID int) (string, *model.User, error) {
	user, err := s.userRepo.Create(ctx, userID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to login or create user: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"issuedAt": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, user, nil
}
