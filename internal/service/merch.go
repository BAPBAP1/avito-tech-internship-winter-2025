package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/repository"
)

var (
	ErrMerchNotFound = errors.New("merch not found")
)

type MerchService struct {
	merchRepo       *repository.MerchRepository
	userRepo        *repository.UserRepository
	transactionRepo *repository.TransactionRepository
	db              *sql.DB
}

func NewMerchService(merchRepo *repository.MerchRepository, userRepo *repository.UserRepository, db *sql.DB) *MerchService {
	return &MerchService{
		merchRepo:       merchRepo,
		userRepo:        userRepo,
		transactionRepo: repository.NewTransactionRepository(db),
		db:              db,
	}
}

func (s *MerchService) ListMerch(ctx context.Context) ([]model.Merch, error) {
	merchItems, err := s.merchRepo.ListMerchItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list merch items: %w", err)
	}
	return merchItems, nil
}

func (s *MerchService) PurchaseMerch(ctx context.Context, userID int, itemName string) error {
	merchItem, err := s.merchRepo.GetMerchItemByName(ctx, itemName)
	if err != nil {
		return ErrMerchNotFound
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	userRepoTx := repository.NewUserRepositoryWithTx(tx)
	user, err := userRepoTx.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Coins < merchItem.Price {
		return ErrInsufficientFunds
	}

	newBalance := user.Coins - merchItem.Price
	if err = userRepoTx.UpdateCoins(ctx, userID, newBalance); err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}

	transactionRepoTx := repository.NewTransactionRepositoryWithTx(tx)
	if err = transactionRepoTx.CreatePurchase(ctx, userID, itemName, merchItem.Price); err != nil {
		return fmt.Errorf("failed to record purchase: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *MerchService) ListPurchases(ctx context.Context, userID int) ([]model.Purchase, error) {
	purchases, err := s.transactionRepo.GetPurchasesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchases for user %d: %w", userID, err)
	}
	return purchases, nil
}

func (s *MerchService) CreatePurchaseForUser(ctx context.Context, userID int, itemName string) error {
	merchItem, err := s.merchRepo.GetMerchItemByName(ctx, itemName)
	if err != nil {
		return ErrMerchNotFound
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	transactionRepoTx := repository.NewTransactionRepositoryWithTx(tx)
	if err = transactionRepoTx.CreatePurchase(ctx, userID, itemName, merchItem.Price); err != nil {
		return fmt.Errorf("failed to record purchase: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
