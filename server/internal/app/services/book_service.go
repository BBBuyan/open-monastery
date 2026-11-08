package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type BookService struct {
	repo abstractions.BookRepository
}

func NewBookService(
	repo abstractions.BookRepository,
) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) GetBook(
	ctx context.Context,
	id int64,
) (*models.Book, error) {
	return s.repo.FindById(ctx, id)
}

func (s *BookService) GetBooks(
	ctx context.Context,
	opts appmodels.ListOptions,
) (*appmodels.Paginated[*models.Book], error) {
	opts, err := opts.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooks(ctx, opts)
	if err != nil {
		return nil, err
	}
	var cursor *appmodels.Cursor
	hasMore := len(books) > opts.Limit
	if hasMore {
		books = books[:opts.Limit]
		cursor = &appmodels.Cursor{
			ID: books[len(books)-1].Id,
		}
	}
	paginated := appmodels.ToPage(books, hasMore, cursor)

	return paginated, nil
}

func (s *BookService) CreateBook(
	ctx context.Context,
	name string,
	desc string,
) (*models.Book, error) {
	book := &models.Book{
		Name:        name,
		Description: desc,
	}
	err := s.repo.Add(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) GetBookCategories(
	ctx context.Context,
	bookId int64,
) ([]*models.Category, error) {
	categories, err := s.repo.GetCategories(ctx, bookId)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
