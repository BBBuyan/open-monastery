package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingService struct {
	repo     abstractions.ReadingRepository
}

func NewReadingService(
	repo abstractions.ReadingRepository,
) *ReadingService {
	return &ReadingService{
		repo:     repo,
	}
}

func (s *ReadingService) GetReading(
	ctx context.Context, 
	id int64,
) (*models.Reading, error) {
	return s.repo.FindById(ctx, id)
}

func (s *ReadingService) CreateReading(
	ctx context.Context,
	bookId int64,
	customerId int64,
) (*models.Reading, error) {
	_, err := s.repo.UserExists(ctx, customerId)
	if err != nil {
		return nil, apperr.NewNotFound("user not found", err)
	}
	_, err = s.repo.BookExists(ctx, bookId)
	if err != nil {
		return nil, apperr.NewNotFound("book not found", err)
	}

	reading := models.NewReading(bookId, customerId)
	if err := s.repo.Add(ctx, reading); err != nil {
		return nil, err
	}

	return reading, nil
}
