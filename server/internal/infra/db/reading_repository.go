package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
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
) (*models.Reading, error) {
	query := `SELECT id, book_id, monk_id, customer_id FROM reading WHERE id = ?`

	var foundReading models.Reading

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundReading.Id,
		&foundReading.BookId,
		&foundReading.MonkId,
		&foundReading.CustomerId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("reading not found", err)
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
		return translateWriteError("add reading", err)
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
        "SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)",
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
