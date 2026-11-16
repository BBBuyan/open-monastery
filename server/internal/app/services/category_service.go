package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type CategoryService struct {
	repo abstractions.CategoryRepository
}

func NewCategoryService(repo abstractions.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategory(
	ctx context.Context,
	monasteryID, id int64,
) (*models.Category, error) {
	return s.repo.FindById(ctx, monasteryID, id)
}

func (s *CategoryService) GetCategories(
	ctx context.Context,
	monasteryID int64,
) ([]*models.Category, error) {
	return s.repo.GetAll(ctx, monasteryID)
}

func (s *CategoryService) GetBooksByCategoryId(
	ctx context.Context,
	monasteryID, categoryId int64,
	opts appmodels.ListOptions,
) (*appmodels.Paginated[*models.Book], error) {
	opts, err := opts.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooksByCategoryId(ctx, monasteryID, categoryId, opts)
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

func (s *CategoryService) AddBook(
	ctx context.Context,
	monasteryID, categoryId, bookId int64,
) error {
	return s.repo.AddBook(ctx, monasteryID, categoryId, bookId)
}

func (s *CategoryService) RemoveBook(
	ctx context.Context,
	monasteryID, categoryId, bookId int64,
) error {
	return s.repo.RemoveBook(ctx, monasteryID, categoryId, bookId)
}

func (s *CategoryService) CreateCategory(
	ctx context.Context,
	monasteryID int64,
	name string,
) (*models.Category, error) {
	category := &models.Category{Name: name, MonasteryId: monasteryID}
	if err := s.repo.Add(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}
