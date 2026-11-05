package category

import (
	"context"

	"example.com/open-monastery/internal/book"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*Category, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Service) GetCategories(ctx context.Context) ([]Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetBooksByCategoryId(ctx context.Context, categoryId int64, options book.ListOptions) (*book.Page, error) {
	options, err := options.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooksByCategoryId(ctx, categoryId, options)
	if err != nil {
		return nil, err
	}
	return book.NewPage(books, options.Limit), nil
}

func (s *Service) GetCategoriesByBookId(ctx context.Context, bookId int64) ([]Category, error) {
	return s.repo.GetCategoriesByBookId(ctx, bookId)
}

func (s *Service) GetCategoryBooks(ctx context.Context) ([]CategoryBook, error) {
	return s.repo.GetCategoryBooks(ctx)
}

func (s *Service) AddBook(ctx context.Context, categoryId, bookId int64) error {
	return s.repo.AddBook(ctx, categoryId, bookId)
}

func (s *Service) RemoveBook(ctx context.Context, categoryId, bookId int64) error {
	return s.repo.RemoveBook(ctx, categoryId, bookId)
}

func (s *Service) CreateCategory(ctx context.Context, name string) (*Category, error) {
	category := &Category{Name: name}
	if err := s.repo.Add(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}
