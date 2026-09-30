package reading

import (
	"context"

	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/user"
)

type Service struct {
	repo     Repository
	userRepo user.Repository
	bookRepo book.Repository
}

func NewService(
	repo Repository,
	userRepo user.Repository,
	bookRepo book.Repository,
) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		bookRepo: bookRepo,
	}
}

func (s *Service) GetReading(ctx context.Context, id int64) (*Reading, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Service) CreateReading(
	ctx context.Context,
	bookId int64,
	customerId int64,
) (*Reading, error) {
	_, err := s.userRepo.FindById(ctx, customerId)
	if err != nil {
		return nil, err
	}
	_, err = s.bookRepo.FindById(ctx, bookId)
	if err != nil {
		return nil, err
	}

	reading := NewReading(bookId, customerId)
	if err := s.repo.Add(ctx, reading); err != nil {
		return nil, err
	}

	return reading, nil
}
