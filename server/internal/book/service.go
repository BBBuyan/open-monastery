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

func (s *Service) CreateBook(
	ctx context.Context,
	name string,
	description string,
) (*Book, error) {
	book := &Book{
		Name: name,
		Description: description,
	}
	err := s.repo.Add(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}
