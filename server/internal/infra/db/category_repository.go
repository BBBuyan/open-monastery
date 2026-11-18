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
	monasteryID, categoryId int64,
	opts appmodels.ListOptions,
) ([]*models.Book, error) {
	query := `
	SELECT b.id, b.name, b.description, b.price, b.monastery_id
	FROM book b
		WHERE b.monastery_id = ? AND b.id > ? AND EXISTS (
			SELECT 1 FROM category_book cb
			JOIN category c ON c.id = cb.category_id AND c.monastery_id = b.monastery_id
			WHERE cb.book_id = b.id AND cb.category_id = ?
		)
	`

	args := []any{monasteryID, opts.AfterId, categoryId}
	if opts.Search != "" {
		query += ` AND (LOCATE(?, b.name) > 0 OR LOCATE(?, b.description) > 0)`
		args = append(args, opts.Search, opts.Search)
	}
	query += " ORDER BY b.id LIMIT ?"
	args = append(args, opts.Limit+1)
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
			&found.MonasteryId,
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
	monasteryID int64,
) ([]*models.CategoryBook, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT cb.id, cb.category_id, cb.book_id FROM category_book cb
		JOIN category c ON c.id = cb.category_id
		JOIN book b ON b.id = cb.book_id AND b.monastery_id = c.monastery_id
		WHERE c.monastery_id = ? ORDER BY cb.id`,
		monasteryID,
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
	monasteryID, categoryId, bookId int64,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO category_book (category_id, book_id)
		SELECT c.id, b.id FROM category c
		JOIN book b ON b.id = ? AND b.monastery_id = c.monastery_id
		WHERE c.id = ? AND c.monastery_id = ?`,
		bookId,
		categoryId,
		monasteryID,
	)
	if err != nil {
		return translateWriteError("AddBook", err)
	}
	added, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get number of books added to category: %w", err)
	}
	if added == 0 {
		return apperr.NewNotFound("category or book not found", nil)
	}
	return nil
}

func (r *CategoryRepository) RemoveBook(
	ctx context.Context,
	monasteryID, categoryId, bookId int64,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE cb FROM category_book cb
		WHERE cb.category_id = ? AND cb.book_id = ? AND cb.monastery_id = ?`,
		categoryId,
		bookId,
		monasteryID,
	)
	if err != nil {
		return translateWriteError("RemoveBook", err)
	}
	return nil
}

func (r *CategoryRepository) FindById(
	ctx context.Context,
	monasteryID, id int64,
) (*models.Category, error) {
	var found models.Category
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, monastery_id FROM category WHERE monastery_id = ? AND id = ?`,
		monasteryID,
		id,
	).Scan(&found.Id, &found.Name, &found.MonasteryId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.NewNotFound(
			"category not found",
			err,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}
	return &found, nil
}

func (r *CategoryRepository) GetAll(
	ctx context.Context,
	monasteryID int64,
) ([]*models.Category, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, monastery_id FROM category WHERE monastery_id = ? ORDER BY id`,
		monasteryID,
	)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var found models.Category
		if err := rows.Scan(&found.Id, &found.Name, &found.MonasteryId); err != nil {
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
		`INSERT INTO category (name, monastery_id) VALUES (?, ?)`,
		category.Name,
		category.MonasteryId,
	)
	if err != nil {
		return translateWriteError("AddCategory", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted category ID: %w", err)
	}
	category.Id = id
	return nil
}
