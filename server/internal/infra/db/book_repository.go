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

type BookRepository struct {
	db *sql.DB
}

var _ abstractions.BookRepository = (*BookRepository)(nil)

func (r *BookRepository) GetBooks(
	ctx context.Context,
	options appmodels.ListOptions,
) ([]*models.Book, error) {
	query := `
		SELECT id, name, description, price 
		FROM book WHERE id > ?
	`
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
	books := make([]*models.Book, 0)
	for rows.Next() {
		var found models.Book
		if err := rows.Scan(
			&found.Id,
			&found.Name,
			&found.Description,
			&found.Price,
		); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, &found)
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
) (*models.Book, error) {
	query := `SELECT id, name, description, price FROM book WHERE id = ?`

	var foundBook models.Book
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundBook.Id,
		&foundBook.Name,
		&foundBook.Description,
		&foundBook.Price,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("book not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("find book: %w", err)
	}

	return &foundBook, nil
}

func (r *BookRepository) Add(
	ctx context.Context,
	book *models.Book,
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

func (r *BookRepository) GetCategories(
	ctx context.Context,
	bookId int64,
) ([]*models.Category, error) {
	query := `SELECT c.id, c.name FROM category c
		WHERE EXISTS (
			SELECT 1 FROM category_book cb 
			WHERE cb.category_id = c.id AND cb.book_id = ?
		) ORDER BY c.id`
	rows, err := r.db.QueryContext(ctx, query, bookId)
	if err != nil {
		return nil, fmt.Errorf("get categories by book ID: %w", err)
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var found models.Category
		if err := rows.Scan(&found.Id, &found.Name); err != nil {
			return nil, fmt.Errorf("scan book category: %w", err)
		}
		categories = append(categories, &found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories by book ID: %w", err)
	}
	return categories, nil
}
