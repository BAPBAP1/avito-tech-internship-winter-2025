package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func NewTransactionRepositoryWithTx(tx *sql.Tx) *TransactionRepository {
	return &TransactionRepository{tx: tx}
}

func (r *TransactionRepository) Create(ctx context.Context, senderID, receiverID int, amount int) error {
	var execContext func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	if r.tx != nil {
		execContext = r.tx.ExecContext
	} else {
		execContext = r.db.ExecContext
	}

	_, err := execContext(ctx,
		"INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)",
		senderID, receiverID, amount,
	)
	return err
}

func (r *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID int) ([]model.Transaction, error) {
	var queryContext func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	if r.tx != nil {
		queryContext = r.tx.QueryContext
	} else {
		queryContext = r.db.QueryContext
	}

	rows, err := queryContext(ctx,
		`SELECT id, sender_id, receiver_id, amount, created_at
   FROM transactions
   WHERE sender_id = $1 OR receiver_id = $1
   ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Transaction{}, nil 
		}
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) CreatePurchase(ctx context.Context, userID int, itemName string, price int) error {
	var execContext func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	if r.tx != nil {
		execContext = r.tx.ExecContext
	} else {
		execContext = r.db.ExecContext
	}

	_, err := execContext(ctx,
		"INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)",
		userID, itemName, price,
	)
	return err
}

func (r *TransactionRepository) GetPurchasesByUserID(ctx context.Context, userID int) ([]model.Purchase, error) {
	var queryContext func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	if r.tx != nil {
		queryContext = r.tx.QueryContext
	} else {
		queryContext = r.db.QueryContext
	}

	rows, err := queryContext(ctx,
		`SELECT id, user_id, item_name, price, purchased_at
   FROM purchases
   WHERE user_id = $1
   ORDER BY purchased_at DESC`, userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Purchase{}, nil 
		}
		return nil, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []model.Purchase
	for rows.Next() {
		var p model.Purchase
		var purchasedAt time.Time
		if err := rows.Scan(&p.ID, &p.UserID, &p.ItemName, &p.Price, &purchasedAt); err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}
		p.PurchasedAt = purchasedAt.Format(time.RFC3339)
		purchases = append(purchases, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating purchase rows: %w", err)
	}

	return purchases, nil
}
