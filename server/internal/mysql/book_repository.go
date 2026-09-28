package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/book"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) FindById(
	ctx context.Context,
	id int64,
) (*book.Book, error) {
	query := `SELECT id, name, description, price FROM book WHERE id = ?`

	var foundBook book.Book
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundBook.Id,
		&foundBook.Name,
		&foundBook.Description,
		&foundBook.Price,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.NotFound, "book not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("find book: %w", err)
	}

	return &foundBook, nil
}

func (r *BookRepository) Add(
	ctx context.Context,
	book *book.Book,
) error {
	query := `INSERT INTO book (name,description,price) VALUES(?,?,?)`
	result, err := r.db.ExecContext(
		ctx,
		query,
		book.Name,
		book.Description,
		book.Price,
	)

	if err != nil {
		return translateWriteError("add book", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted book ID: %w", err)
	}
	book.Id = id

	return nil
}
