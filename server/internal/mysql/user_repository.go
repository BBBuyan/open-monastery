package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/user"
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
) (*user.User, error) {
	query := `SELECT id, first_name, last_name, user_type, monastery_id FROM users WHERE id = ?`
	var foundUser user.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundUser.Id,
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.UserType,
		&foundUser.MonasteryId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.NotFound, "user not found", err)
	}

	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return &foundUser, nil
}

func (r *UserRepository) Add(
	ctx context.Context,
	user *user.User,
) error {
	query := `INSERT INTO users (first_name, last_name, user_type, monastery_id) VALUES (?,?,?,?)`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserType,
		user.MonasteryId,
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
