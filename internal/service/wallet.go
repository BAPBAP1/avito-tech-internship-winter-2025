package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/repository"
)

var (
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type WalletService struct {
	userRepo        *repository.UserRepository
	transactionRepo *repository.TransactionRepository
	db              *sql.DB
}

func NewWalletService(userRepo *repository.UserRepository, transactionRepo *repository.TransactionRepository, db *sql.DB) *WalletService {
	return &WalletService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *WalletService) Transfer(ctx context.Context, senderID, receiverID int, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
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

	// Используем репозитории с поддержкой транзакций
	userRepoTx := repository.NewUserRepositoryWithTx(tx)
	transactionRepoTx := repository.NewTransactionRepositoryWithTx(tx)

	// Получаем отправителя и проверяем его баланс
	sender, err := userRepoTx.GetByID(ctx, senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender user: %w", err)
	}

	if sender.Coins < amount {
		return ErrInsufficientFunds
	}

	// Получаем получателя
	receiver, err := userRepoTx.GetByID(ctx, receiverID)
	if err != nil {
		return fmt.Errorf("failed to get receiver user: %w", err)
	}

	// Обновляем балансы
	senderNewBalance := sender.Coins - amount
	receiverNewBalance := receiver.Coins + amount

	if err := userRepoTx.UpdateCoins(ctx, senderID, senderNewBalance); err != nil {
		return fmt.Errorf("failed to update sender coins: %w", err)
	}

	if err := userRepoTx.UpdateCoins(ctx, receiverID, receiverNewBalance); err != nil {
		return fmt.Errorf("failed to update receiver coins: %w", err)
	}

	// Записываем транзакцию
	if err := transactionRepoTx.Create(ctx, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *WalletService) GetWallet(ctx context.Context, userID int) (*model.Wallet, error) {
	coins, err := s.userRepo.GetCoins(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get coins for user %d: %w", userID, err)
	}

	history, err := s.GetWalletHistory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet history for user %d: %w", userID, err)
	}

	return &model.Wallet{
		Coins:              coins,
		TransactionHistory: history,
	}, nil
}

func (s *WalletService) GetWalletHistory(ctx context.Context, userID int) ([]model.WalletHistoryEntry, error) {
	transactions, err := s.transactionRepo.GetTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for user %d: %w", userID, err)
	}

	var historyEntries []model.WalletHistoryEntry
	for _, tx := range transactions {
		entry := model.WalletHistoryEntry{
			Amount:          tx.Amount,
			CreatedAt:       tx.CreatedAt.Format(time.RFC3339),
			CounterpartyID:  tx.SenderID,
			TransactionType: "incoming",
		}

		if tx.SenderID == userID {
			entry.TransactionType = "outgoing"
			entry.CounterpartyID = tx.ReceiverID
		}

		historyEntries = append(historyEntries, entry)
	}

	return historyEntries, nil
}