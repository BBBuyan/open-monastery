package book

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetBook(ctx context.Context, id int64) (*Book, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Service) GetBooks(ctx context.Context, options ListOptions) (*Page, error) {
	options, err := options.Normalize()
	if err != nil {
		return nil, err
	}
	books, err := s.repo.GetBooks(ctx, options)
	if err != nil {
		return nil, err
	}
	return NewPage(books, options.Limit), nil
}

func (s *Service) CreateBook(
	ctx context.Context,
	name string,
	description string,
) (*Book, error) {
	book := &Book{
		Name:        name,
		Description: description,
	}
	err := s.repo.Add(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}
