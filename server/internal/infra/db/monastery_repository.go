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
