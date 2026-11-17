package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type MonasteryRepository struct {
	db *sql.DB
}

func (r *MonasteryRepository) CanWrite(ctx context.Context, monasteryID, userID int64) (bool, error) {
	var allowed bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE id = ?
			AND (user_type = ? OR (user_type = ? AND monastery_id = ?))
		)
	`, userID, models.Admin, models.Manager, monasteryID).Scan(&allowed)
	if err != nil {
		return false, fmt.Errorf("check monastery write access: %w", err)
	}
	return allowed, nil
}

func (r *MonasteryRepository) GrantManager(ctx context.Context, monasteryID, userID int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE users SET user_type = ?, monastery_id = ?
		WHERE id = ? AND user_type = ?
	`, models.Manager, monasteryID, userID, models.Customer)
	if err != nil {
		return translateWriteError("grant monastery manager", err)
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count appointed monastery managers: %w", err)
	}
	if updated > 0 {
		return nil
	}

	var role models.UserType
	var assignedMonastery sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT user_type, monastery_id FROM users WHERE id = ?
	`, userID).Scan(&role, &assignedMonastery)
	if errors.Is(err, sql.ErrNoRows) {
		return apperr.NewNotFound("user.not_found", err)
	}
	if err != nil {
		return fmt.Errorf("check existing manager appointment: %w", err)
	}
	if role == models.Manager && assignedMonastery.Valid && assignedMonastery.Int64 == monasteryID {
		return nil
	}
	return apperr.NewConflict("user.not_eligible_for_manager", nil)
}

func (r *MonasteryRepository) RevokeManager(ctx context.Context, monasteryID, userID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET user_type = ?, monastery_id = NULL
		WHERE id = ? AND monastery_id = ? AND user_type = ?
	`, models.Customer, userID, monasteryID, models.Manager)
	if err != nil {
		return fmt.Errorf("revoke monastery manager: %w", err)
	}
	return nil
}

var _ abstractions.MonasteryRepository = (*MonasteryRepository)(nil)

func NewMonasteryRepo(db *sql.DB) *MonasteryRepository {
	return &MonasteryRepository{
		db: db,
	}
}

func (r *MonasteryRepository) FindById(
	ctx context.Context,
	id int64,
) (*models.Monastery, error) {
	query := `SELECT id, name FROM monastery WHERE id = ?`

	var foundMonastery models.Monastery

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundMonastery.Id,
		&foundMonastery.Name,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound(
			"monastery not found",
			err,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("find monastery: %w", err)
	}
	return &foundMonastery, nil
}

func (r *MonasteryRepository) Add(
	ctx context.Context,
	monastery *models.Monastery,
) error {
	query := `INSERT INTO monastery (name) VALUES(?)`
	result, err := r.db.ExecContext(ctx, query, monastery.Name)

	if err != nil {
		return translateWriteError("AddMonastery", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted monastery ID: %w", err)
	}
	monastery.Id = id

	return nil
}

func (r *MonasteryRepository) GetMonasteries(
	ctx context.Context,
	opts appmodels.ListOptions,
) ([]*models.Monastery, error) {
	query := "SELECT id, name FROM monastery WHERE id > ?"
	args := []any{opts.AfterId}
	if opts.Search != "" {
		query += " AND LOCATE(?, name) > 0"
		args = append(args, opts.Search)
	}
	query += " ORDER BY id LIMIT ?"
	args = append(args, opts.Limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get monasteries: %w", err)
	}
	defer rows.Close()

	monasteries := make([]*models.Monastery, 0)
	for rows.Next() {
		var monastery models.Monastery
		if err := rows.Scan(&monastery.Id, &monastery.Name); err != nil {
			return nil, fmt.Errorf("scan monastery: %w", err)
		}
		monasteries = append(monasteries, &monastery)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monasteries: %w", err)
	}
	return monasteries, nil
}
