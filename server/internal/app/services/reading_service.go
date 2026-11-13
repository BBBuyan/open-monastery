package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingService struct {
	repo abstractions.ReadingRepository
}

func NewReadingService(
	repo abstractions.ReadingRepository,
) *ReadingService {
	return &ReadingService{
		repo: repo,
	}
}

func (s *ReadingService) GetReadings(
	ctx context.Context,
	userID int64,
	opts appmodels.ListOptions,
) (*appmodels.Paginated[*models.Reading], error) {
	opts, err := opts.Normalize()
	if err != nil {
		return nil, err
	}
	readings, err := s.repo.GetReadingsByUserId(
		ctx,
		userID,
		opts,
	)
	if err != nil {
		return nil, err
	}
	var cursor *appmodels.Cursor
	hasMore := len(readings) > opts.Limit
	if hasMore {
		readings = readings[:opts.Limit]
		cursor = &appmodels.Cursor{
			ID: readings[len(readings)-1].Id,
		}
	}
	paginated := appmodels.ToPage(readings, hasMore, cursor)

	return paginated, nil
}

func (s *ReadingService) CreateReading(
	ctx context.Context,
	bookId int64,
	customerId int64,
) (*models.Reading, error) {
	userExists, err := s.repo.UserExists(ctx, customerId)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, apperr.NewNotFound("user not found", nil)
	}
	bookExists, err := s.repo.BookExists(ctx, bookId)
	if err != nil {
		return nil, err
	}
	if !bookExists {
		return nil, apperr.NewNotFound("book not found", nil)
	}

	reading := models.NewReading(bookId, customerId)
	if err := s.repo.Add(ctx, reading); err != nil {
		return nil, err
	}

	return reading, nil
}
