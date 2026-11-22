package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
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
	monasteryId, id int64,
) (*models.Category, error) {
	return s.repo.FindById(ctx, monasteryId, id)
}

func (s *CategoryService) GetCategories(
	ctx context.Context,
	monasteryId int64,
) ([]*models.Category, error) {
	return s.repo.GetAll(ctx, monasteryId)
}

func (s *CategoryService) GetBooksByCategoryId(
	ctx context.Context,
	monasteryId, categoryId int64,
	opts appmodels.ListOptions,
) (*appmodels.Paginated[*models.Book], error) {
	opts, err := opts.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooksByCategoryId(ctx, monasteryId, categoryId, opts)
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
	categoryId, bookId int64,
	actor appmodels.Actor,
) error {
	if actor.Role != enums.Manager {
		return apperr.NewForbidden("forbidden")
	}
	if actor.MonasteryId == nil {
		return apperr.NewBadRequest("user has no monastery id", nil)
	}
	categoryBook := models.NewCategoryBook(categoryId, bookId, *actor.MonasteryId, &actor.Email)

	return s.repo.AddCategoryBook(ctx, categoryBook)
}

func (s *CategoryService) RemoveBook(
	ctx context.Context,
	categoryId, bookId int64,
	actor appmodels.Actor,
) error {
	if actor.Role != enums.Manager {
		return apperr.NewForbidden("forbidden")
	}
	if actor.MonasteryId == nil {
		return apperr.NewBadRequest("user has no monastery id", nil)
	}
	return s.repo.RemoveBook(ctx, *actor.MonasteryId, categoryId, bookId)
}

func (s *CategoryService) CreateCategory(
	ctx context.Context,
	name string,
	actor appmodels.Actor,
) (*models.Category, error) {
	if actor.Role != enums.Manager {
		return nil, apperr.NewForbidden("forbidden")
	}
	if actor.MonasteryId == nil {
		return nil, apperr.NewBadRequest("user has no monastery id", nil)
	}
	category := models.NewCategory(name, *actor.MonasteryId, &actor.Email)
	if err := s.repo.Add(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}
