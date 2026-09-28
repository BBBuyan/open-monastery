package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/monastery"
)

type MonasteryRepository struct {
	db *sql.DB
}

func NewMonasteryRepo(db *sql.DB) *MonasteryRepository {
	return &MonasteryRepository{
		db: db,
	}
}

func (r *MonasteryRepository) FindById(
	ctx context.Context,
	id int64,
) (*monastery.Monastery, error) {
	query := `SELECT id, name FROM monastery WHERE id = ?`

	var foundMonastery monastery.Monastery

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundMonastery.Id,
		&foundMonastery.Name,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.NotFound, "monastery not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("find monastery: %w", err)
	}
	return &foundMonastery, nil
}

func (r *MonasteryRepository) Add(
	ctx context.Context,
	monastery *monastery.Monastery,
) error {
	query := `INSERT INTO monastery (name) VALUES(?)`
	result, err := r.db.ExecContext(ctx, query, monastery.Name)

	if err != nil {
		return translateWriteError("add monastery", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted monastery ID: %w", err)
	}
	monastery.Id = id

	return nil
}
