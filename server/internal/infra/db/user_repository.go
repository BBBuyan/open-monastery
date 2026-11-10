package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindById(
	ctx context.Context,
	id int64,
) (*models.User, error) {
	query := `SELECT id, first_name, last_name, user_type, monastery_id, email, phone_number
	FROM users WHERE id = ?`
	var foundUser models.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundUser.Id,
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.UserType,
		&foundUser.MonasteryId,
		&foundUser.Email,
		&foundUser.PhoneNumber,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound(
			"user not found",
			"user not found",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return &foundUser, nil
}

func (r *UserRepository) Add(
	ctx context.Context,
	user *models.User,
) error {
	query := `INSERT INTO users (first_name, last_name, user_type, monastery_id, email, phone_number)
	VALUES (?,?,?,?,?,?)`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserType,
		user.MonasteryId,
		user.Email,
		user.PhoneNumber,
	)
	if err != nil {
		return translateWriteError("add user", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted user ID: %w", err)
	}
	user.Id = id

	return nil
}
