package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/commands"
	"example.com/open-monastery/internal/domain/enums"
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
	monasteryID, id int64,
) (*models.Book, error) {
	return s.repo.FindById(ctx, monasteryID, id)
}

func (s *BookService) GetBooks(
	ctx context.Context,
	monasteryID int64,
	opts appmodels.ListOptions,
) (*appmodels.Paginated[*models.Book], error) {
	opts, err := opts.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooks(ctx, monasteryID, opts)
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
	cmd commands.CreateBookCommand,
	actor appmodels.Actor,
) (*models.Book, error) {
	if actor.Role != enums.Manager {
		return nil, apperr.NewForbidden("forbidden")
	}
	if actor.MonasteryId == nil {
		return nil, apperr.NewBadRequest("user has no monastery id", nil)
	}
	book := models.NewBook(
		cmd.Name,
		cmd.Desc,
		cmd.Price,
		*actor.MonasteryId,
		actor.Id,
	)
	err := s.repo.Add(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) GetBookCategories(
	ctx context.Context,
	monasteryID, bookId int64,
) ([]*models.Category, error) {
	categories, err := s.repo.GetCategories(ctx, monasteryID, bookId)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
