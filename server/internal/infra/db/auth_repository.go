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

type AuthRepository struct {
	db *sql.DB
}

var _ abstractions.AuthRepository = (*AuthRepository)(nil)

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindCredentialsByEmail(
	ctx context.Context,
	email string,
) (int64, string, models.UserType, error) {
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

func (r *AuthRepository) AddUser(
	ctx context.Context,
	user *models.User,
	passwordHash string,
) error {
	query := `INSERT INTO users (first_name, last_name, user_type, monastery_id, email, phone_number, password_hash)
	VALUES (?,?,?,?,?,?,?)`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserType,
		user.MonasteryId,
		user.Email,
		user.PhoneNumber,
		passwordHash,
	)
	if err != nil {
		return translateWriteError("AddUser", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted user ID: %w", err)
	}
	user.Id = id

	return nil
}
