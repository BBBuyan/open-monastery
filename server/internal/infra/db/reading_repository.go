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

type ReadingRepository struct {
	db *sql.DB
}

var _ abstractions.ReadingRepository = (*ReadingRepository)(nil)

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	return &ReadingRepository{
		db: db,
	}
}

func (r *ReadingRepository) FindByIdForCustomer(
	ctx context.Context,
	id int64,
	customerID int64,
) (*models.Reading, error) {
	query := `SELECT id, book_id, monk_id, customer_id FROM reading WHERE id = ? AND customer_id = ?`

	var foundReading models.Reading

	err := r.db.QueryRowContext(ctx, query, id, customerID).Scan(
		&foundReading.Id,
		&foundReading.BookId,
		&foundReading.MonkId,
		&foundReading.CustomerId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound(
			"reading not found",
			err,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("find reading: %w", err)
	}

	return &foundReading, nil
}

func (r *ReadingRepository) Add(
	ctx context.Context,
	reading *models.Reading,
) error {
	query := `INSERT INTO reading (book_id, monk_id, customer_id) VALUES(?,?,?)`
	result, err := r.db.ExecContext(
		ctx,
		query,
		reading.BookId,
		reading.MonkId,
		reading.CustomerId,
	)

	if err != nil {
		return translateWriteError("AddReading", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted reading ID: %w", err)
	}
	reading.Id = id

	return nil
}

func (r *ReadingRepository) UserExists(
	ctx context.Context,
	id int64,
) (exists bool, err error) {
	err = r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		id,
	).Scan(&exists)

	return
}

func (r *ReadingRepository) BookExists(
	ctx context.Context,
	id int64,
) (exists bool, err error) {
	err = r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM book WHERE id = ?)",
		id,
	).Scan(&exists)

	return
}

// GetReadingsByUserId implements [abstractions.ReadingRepository].
func (r *ReadingRepository) GetReadingsByUserId(
	ctx context.Context,
	id int64,
	opts appmodels.ListOptions,
) ([]*models.Reading, error) {
	query := `
	SELECT id, book_id, monk_id, customer_id
	FROM reading 
	WHERE user_id = ? AND id > ?
	ORDER BY id LIMIT ?
	`
	args := []any{id, opts.AfterId, opts.Limit + 1}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get readings by user_id: %w", err)
	}
	defer rows.Close()
	readings := make([]*models.Reading, 0)
	for rows.Next() {
		var found models.Reading
		if err := rows.Scan(
			&found.Id,
			&found.BookId,
			&found.MonkId,
			&found.CustomerId,
		); err != nil {
			return nil, fmt.Errorf("scan reading: %w", err)
		}
		readings = append(readings, &found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate readings by user id: %w", err)
	}

	return readings, nil
}
