package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
)

type LoginRepository struct {
	db *sql.DB
}

var _ abstractions.LoginRepository = (*LoginRepository)(nil)

func NewLoginRepository(db *sql.DB) *LoginRepository {
	return &LoginRepository{db: db}
}

func (r *LoginRepository) FindPasswordHash(ctx context.Context, userID int64) (string, error) {
	const query = `SELECT password_hash FROM users WHERE id = ? AND password_hash IS NOT NULL`
	var hash string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperr.NewNotFound("credentials not found", err)
	}
	if err != nil {
		return "", fmt.Errorf("find password hash: %w", err)
	}
	return hash, nil
}
