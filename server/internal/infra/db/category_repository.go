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

type CategoryRepository struct {
	db *sql.DB
}

var _ abstractions.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepo(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetBooksByCategoryId(
	ctx context.Context,
	categoryId int64,
	options appmodels.ListOptions,
) ([]*models.Book, error) {
	query := `SELECT b.id, b.name, b.description, b.price 
	FROM book b
		WHERE b.id > ? AND EXISTS (
			SELECT 1 FROM category_book cb 
			WHERE cb.book_id = b.id AND cb.category_id = ?
		)`

	args := []any{options.AfterId, categoryId}
	if options.Search != "" {
		query += ` AND (LOCATE(?, b.name) > 0 OR LOCATE(?, b.description) > 0)`
		args = append(args, options.Search, options.Search)
	}
	query += " ORDER BY b.id LIMIT ?"
	args = append(args, options.Limit+1)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get books by category ID: %w", err)
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
			return nil, fmt.Errorf("scan category book: %w", err)
		}
		books = append(books, &found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books by category ID: %w", err)
	}
	return books, nil
}

func (r *CategoryRepository) GetCategoryBooks(
	ctx context.Context,
) ([]*models.CategoryBook, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, category_id, book_id FROM category_book ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("get category books: %w", err)
	}
	defer rows.Close()

	links := make([]*models.CategoryBook, 0)
	for rows.Next() {
		var link models.CategoryBook
		if err := rows.Scan(
			&link.Id,
			&link.CategoryId,
			&link.BookId,
		); err != nil {
			return nil, fmt.Errorf("scan category book: %w", err)
		}
		links = append(links, &link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate category books: %w", err)
	}
	return links, nil
}

func (r *CategoryRepository) AddBook(
	ctx context.Context,
	categoryId,
	bookId int64,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO category_book (category_id, book_id) VALUES (?, ?)`,
		categoryId,
		bookId,
	)
	if err != nil {
		return translateWriteError("add book to category", err)
	}
	return nil
}

func (r *CategoryRepository) RemoveBook(
	ctx context.Context,
	categoryId int64,
	bookId int64,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM category_book WHERE category_id = ? AND book_id = ?`,
		categoryId,
		bookId,
	)
	if err != nil {
		return translateWriteError("remove book from category", err)
	}
	return nil
}

func (r *CategoryRepository) FindById(
	ctx context.Context,
	id int64,
) (*models.Category, error) {
	var found models.Category
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name FROM category WHERE id = ?`,
		id,
	).Scan(&found.Id, &found.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound("category not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}
	return &found, nil
}

func (r *CategoryRepository) GetAll(
	ctx context.Context,
) ([]*models.Category, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name FROM category ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var found models.Category
		if err := rows.Scan(&found.Id, &found.Name); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, &found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepository) Add(
	ctx context.Context,
	category *models.Category,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO category (name) VALUES (?)`,
		category.Name,
	)
	if err != nil {
		return translateWriteError("add category", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted category ID: %w", err)
	}
	category.Id = id
	return nil
}
