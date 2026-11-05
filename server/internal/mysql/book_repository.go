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

var _ book.Repository = (*BookRepository)(nil)

func (r *BookRepository) GetBooks(ctx context.Context, options book.ListOptions) ([]book.Book, error) {
	query := `SELECT id, name, description, price FROM book WHERE id > ?`
	args := []any{options.AfterId}
	if options.Search != "" {
		query += ` AND (LOCATE(?, name) > 0 OR LOCATE(?, description) > 0)`
		args = append(args, options.Search, options.Search)
	}
	query += " ORDER BY id LIMIT ?"
	args = append(args, options.Limit+1)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get books: %w", err)
	}
	defer rows.Close()
	books := make([]book.Book, 0)
	for rows.Next() {
		var found book.Book
		if err := rows.Scan(&found.Id, &found.Name, &found.Description, &found.Price); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}
	return books, nil
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
