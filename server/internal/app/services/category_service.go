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
	id int64,
) (*models.Category, error) {
	return s.repo.FindById(ctx, id)
}

func (s *CategoryService) GetCategories(
	ctx context.Context,
) ([]*models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) GetBooksByCategoryId(
	ctx context.Context,
	categoryId int64,
	options appmodels.ListOptions,
) ([]*models.Book, error) {
	options, err := options.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooksByCategoryId(ctx, categoryId, options)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *CategoryService) GetCategoriesByBookId(
	ctx context.Context,
	bookId int64,
) ([]*models.Category, error) {
	return s.repo.GetCategoriesByBookId(ctx, bookId)
}

func (s *CategoryService) GetCategoryBooks(
	ctx context.Context,
) ([]*models.CategoryBook, error) {
	return s.repo.GetCategoryBooks(ctx)
}

func (s *CategoryService) AddBook(
	ctx context.Context,
	categoryId,
	bookId int64,
) error {
	return s.repo.AddBook(ctx, categoryId, bookId)
}

func (s *CategoryService) RemoveBook(
	ctx context.Context,
	categoryId,
	bookId int64,
) error {
	return s.repo.RemoveBook(ctx, categoryId, bookId)
}

func (s *CategoryService) CreateCategory(
	ctx context.Context,
	name string,
) (*models.Category, error) {
	category := &models.Category{Name: name}
	if err := s.repo.Add(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}
