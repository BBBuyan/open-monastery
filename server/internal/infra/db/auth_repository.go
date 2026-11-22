package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
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
) (*appmodels.UserCreds, error) {
	const query = `
		SELECT id, password_hash, user_type, monastery_id, email
		FROM users 
		WHERE email = ? AND password_hash IS NOT NULL
	`
	var creds appmodels.UserCreds
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&creds.UserId,
		&creds.PasswordHash,
		&creds.UserType,
		&creds.MonasteryId,
		&creds.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NewNotFound("credentials not found", err)
		}
		return nil, fmt.Errorf("find credentials by email: %w", err)
	}
	return &creds, nil
}

func (r *AuthRepository) AddUser(
	ctx context.Context,
	user *models.User,
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
		user.Email.String(),
		user.PhoneNumber.String(),
		user.PasswordHash,
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

func (r *AuthRepository) PromoteToMonk(
	ctx context.Context,
	userId, monasteryId int64,
) error {
	query := `
		UPDATE users 
		SET user_type = ?, monastery_id = ? 
		where id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		enums.Monk,
		monasteryId,
		userId,
	)
	return err
}

func (r *AuthRepository) PromoteToManager(
	ctx context.Context,
	userId, monasteryId int64,
) error {
	query := `
		UPDATE users 
		SET user_type = ?, 
		monastery_id = ? 
		where id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		enums.Manager,
		monasteryId,
		userId,
	)

	return err
}
