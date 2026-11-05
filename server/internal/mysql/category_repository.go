package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/category"
)

type CategoryRepository struct {
	db *sql.DB
}

var _ category.Repository = (*CategoryRepository)(nil)

func NewCategoryRepo(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetBooksByCategoryId(ctx context.Context, categoryId int64, options book.ListOptions) ([]book.Book, error) {
	query := `SELECT b.id, b.name, b.description, b.price FROM book b
		WHERE b.id > ? AND EXISTS (
			SELECT 1 FROM category_book cb WHERE cb.book_id = b.id AND cb.category_id = ?
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

	books := make([]book.Book, 0)
	for rows.Next() {
		var found book.Book
		if err := rows.Scan(&found.Id, &found.Name, &found.Description, &found.Price); err != nil {
			return nil, fmt.Errorf("scan category book: %w", err)
		}
		books = append(books, found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books by category ID: %w", err)
	}
	return books, nil
}

func (r *CategoryRepository) GetCategoriesByBookId(ctx context.Context, bookId int64) ([]category.Category, error) {
	query := `SELECT c.id, c.name FROM category c
		WHERE EXISTS (
			SELECT 1 FROM category_book cb WHERE cb.category_id = c.id AND cb.book_id = ?
		) ORDER BY c.id`
	rows, err := r.db.QueryContext(ctx, query, bookId)
	if err != nil {
		return nil, fmt.Errorf("get categories by book ID: %w", err)
	}
	defer rows.Close()

	categories := make([]category.Category, 0)
	for rows.Next() {
		var found category.Category
		if err := rows.Scan(&found.Id, &found.Name); err != nil {
			return nil, fmt.Errorf("scan book category: %w", err)
		}
		categories = append(categories, found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories by book ID: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoryBooks(ctx context.Context) ([]category.CategoryBook, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, category_id, book_id FROM category_book ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("get category books: %w", err)
	}
	defer rows.Close()

	links := make([]category.CategoryBook, 0)
	for rows.Next() {
		var link category.CategoryBook
		if err := rows.Scan(&link.Id, &link.CategoryId, &link.BookId); err != nil {
			return nil, fmt.Errorf("scan category book: %w", err)
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate category books: %w", err)
	}
	return links, nil
}

func (r *CategoryRepository) AddBook(ctx context.Context, categoryId, bookId int64) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO category_book (category_id, book_id) VALUES (?, ?)`, categoryId, bookId)
	if err != nil {
		return translateWriteError("add book to category", err)
	}
	return nil
}

func (r *CategoryRepository) RemoveBook(ctx context.Context, categoryId, bookId int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM category_book WHERE category_id = ? AND book_id = ?`, categoryId, bookId)
	if err != nil {
		return translateWriteError("remove book from category", err)
	}
	return nil
}

func (r *CategoryRepository) FindById(ctx context.Context, id int64) (*category.Category, error) {
	var found category.Category
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM category WHERE id = ?`, id).Scan(&found.Id, &found.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.NotFound, "category not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}
	return &found, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]category.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM category ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	defer rows.Close()

	categories := make([]category.Category, 0)
	for rows.Next() {
		var found category.Category
		if err := rows.Scan(&found.Id, &found.Name); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepository) Add(ctx context.Context, category *category.Category) error {
	result, err := r.db.ExecContext(ctx, `INSERT INTO category (name) VALUES (?)`, category.Name)
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
