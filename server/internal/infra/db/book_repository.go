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
	monasteryID int64,
	options appmodels.ListOptions,
) ([]*models.Book, error) {
	query := `
		SELECT id, name, description, price, monastery_id, created_by, created_at
		FROM book WHERE monastery_id = ? AND id > ?
	`
	args := []any{monasteryID, options.AfterId}
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
			&found.MonasteryId,
			&found.CreatedBy,
			&found.CreatedAt,
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
	monasteryID, id int64,
) (*models.Book, error) {
	query := `SELECT id, name, description, price, monastery_id, created_by, created_at FROM book WHERE monastery_id = ? AND id = ?`

	var foundBook models.Book
	err := r.db.QueryRowContext(ctx, query, monasteryID, id).Scan(
		&foundBook.Id,
		&foundBook.Name,
		&foundBook.Description,
		&foundBook.Price,
		&foundBook.MonasteryId,
		&foundBook.CreatedBy,
		&foundBook.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound(
			"book not found",
			err,
		)
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
	query := `INSERT INTO book (name, description, price, monastery_id, created_by, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(
		ctx,
		query,
		book.Name,
		book.Description,
		book.Price,
		book.MonasteryId,
		book.CreatedBy,
		book.CreatedAt,
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
	monasteryID, bookID int64,
) ([]*models.Category, error) {
	query := `
		SELECT c.id, c.name, c.monastery_id, c.created_by, c.created_at
		FROM category c
		WHERE c.monastery_id = ? AND EXISTS (
			SELECT 1 FROM category_book cb
			JOIN book b ON b.id = cb.book_id AND b.monastery_id = c.monastery_id
			WHERE cb.category_id = c.id AND b.id = ?
		) ORDER BY c.id`
	rows, err := r.db.QueryContext(ctx, query, monasteryID, bookID)
	if err != nil {
		return nil, fmt.Errorf("get categories by book ID: %w", err)
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var found models.Category
		if err := rows.Scan(&found.Id, &found.Name, &found.MonasteryId, &found.CreatedBy, &found.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan book category: %w", err)
		}
		categories = append(categories, &found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories by book ID: %w", err)
	}
	return categories, nil
}
