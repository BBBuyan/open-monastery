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

func (r *LoginRepository) FindCredentialsByEmail(ctx context.Context, email string) (int64, string, error) {
	const query = `SELECT id, password_hash FROM users WHERE email = ? AND password_hash IS NOT NULL`
	var userID int64
	var hash string
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", apperr.NewNotFound("credentials not found", err)
	}
	if err != nil {
		return 0, "", fmt.Errorf("find credentials by email: %w", err)
	}
	return userID, hash, nil
}
