package reading

import (
	"context"

	"example.com/open-monastery/internal/apperror"
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
	monkId int64,
	customerId int64,
	bookId int64,
) (*Reading, error) {
	monk, err := s.userRepo.FindById(ctx, monkId)
	if err != nil {
		return nil, err
	}
	if monk.UserType != user.Monk {
		return nil, apperror.New(
			apperror.Invalid,
			"user is not a monk",
			nil,
		)
	}
	customer, err := s.userRepo.FindById(ctx, customerId)
	if err != nil {
		return nil, err
	}
	if customer.UserType != user.Customer {
		return nil, apperror.New(
			apperror.Invalid,
			"user is not a customer",
			nil,
		)
	}
	foundBook, err := s.bookRepo.FindById(ctx, bookId)
	if err != nil {
		return nil, err
	}
	reading := &Reading{
		CustomerId: customer.Id,
		MonkId:     monk.Id,
		BookId:     foundBook.Id,
	}
	if err := s.repo.Add(ctx, reading); err != nil {
		return nil, err
	}

	return reading, nil
}
