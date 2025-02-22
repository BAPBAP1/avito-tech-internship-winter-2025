package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewUserRepositoryWithTx(tx *sql.Tx) *UserRepository {
	return &UserRepository{tx: tx}
}

func (r *UserRepository) Create(ctx context.Context, userID int) (*model.User, error) {
	var execContext func(ctx context.Context, query string, args ...interface{}) *sql.Row
	if r.tx != nil {
		execContext = r.tx.QueryRowContext
	} else {
		execContext = r.db.QueryRowContext
	}

	var user model.User
	err := execContext(ctx,
		"INSERT INTO users(id, coins) VALUES($1, 1000) ON CONFLICT (id) DO NOTHING RETURNING id, coins", userID,
	).Scan(&user.ID, &user.Coins)

	if err == sql.ErrNoRows {
		return r.GetByID(ctx, userID)
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if user.ID == 0 {
		return r.GetByID(ctx, userID)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	var queryRow func(ctx context.Context, query string, args ...interface{}) *sql.Row
	if r.tx != nil {
		queryRow = r.tx.QueryRowContext
	} else {
		queryRow = r.db.QueryRowContext
	}

	var user model.User
	err := queryRow(ctx,
		"SELECT id, coins FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateCoins(ctx context.Context, id int, newCoins int) error {
	var execContext func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	if r.tx != nil {
		execContext = r.tx.ExecContext
	} else {
		execContext = r.db.ExecContext
	}

	res, err := execContext(ctx,
		"UPDATE users SET coins = $1 WHERE id = $2", newCoins, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows count after update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) GetCoins(ctx context.Context, id int) (int, error) {
	var queryRow func(ctx context.Context, query string, args ...interface{}) *sql.Row
	if r.tx != nil {
		queryRow = r.tx.QueryRowContext
	} else {
		queryRow = r.db.QueryRowContext
	}

	var coins int
	err := queryRow(ctx,
		"SELECT coins FROM users WHERE id = $1", id,
	).Scan(&coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("failed to get user coins: %w", err)
	}
	return coins, nil
}
