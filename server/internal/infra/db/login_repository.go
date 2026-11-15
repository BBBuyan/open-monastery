package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
)

type LoginRepository struct {
	db *sql.DB
}

var _ abstractions.LoginRepository = (*LoginRepository)(nil)

func NewLoginRepository(db *sql.DB) *LoginRepository {
	return &LoginRepository{db: db}
}

func (r *LoginRepository) FindCredentialsByEmail(ctx context.Context, email string) (int64, string, models.UserType, error) {
	const query = `SELECT id, password_hash, user_type FROM users WHERE email = ? AND password_hash IS NOT NULL`
	var userID int64
	var hash string
	var userType models.UserType
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID, &hash, &userType)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", "", apperr.NewNotFound("credentials not found", err)
	}
	if err != nil {
		return 0, "", "", fmt.Errorf("find credentials by email: %w", err)
	}
	return userID, hash, userType, nil
}
