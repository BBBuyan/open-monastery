package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/reading"
)

type ReadingRepository struct {
	db *sql.DB
}

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	return &ReadingRepository{
		db: db,
	}
}

func (r *ReadingRepository) FindById(
	ctx context.Context,
	id int64,
) (*reading.Reading, error) {
	query := `SELECT id, book_id, monk_id, customer_id FROM reading WHERE id = ?`

	var foundReading reading.Reading

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundReading.Id,
		&foundReading.BookId,
		&foundReading.MonkId,
		&foundReading.CustomerId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.NotFound, "reading not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("find reading: %w", err)
	}

	return &foundReading, nil
}

func (r *ReadingRepository) Add(
	ctx context.Context,
	reading *reading.Reading,
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
		return translateWriteError("add reading", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted reading ID: %w", err)
	}
	reading.Id = id

	return nil
}
