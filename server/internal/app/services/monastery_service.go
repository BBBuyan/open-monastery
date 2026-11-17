package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type MonasteryService struct {
	repo abstractions.MonasteryRepository
}

func NewMonasteryService(repo abstractions.MonasteryRepository) *MonasteryService {
	return &MonasteryService{
		repo: repo,
	}
}

func (s *MonasteryService) GetMonastery(
	ctx context.Context,
	id int64,
) (*models.Monastery, error) {
	return s.repo.FindById(ctx, id)
}

func (s *MonasteryService) GetMonasteries(
	ctx context.Context,
	opts appmodels.ListOptions,
) (*appmodels.Paginated[*models.Monastery], error) {
	opts, err := opts.Normalize()
	if err != nil {
		return nil, err
	}
	monasteries, err := s.repo.GetMonasteries(ctx, opts)
	if err != nil {
		return nil, err
	}
	var cursor *appmodels.Cursor
	hasMore := len(monasteries) > opts.Limit
	if hasMore {
		monasteries = monasteries[:opts.Limit]
		cursor = &appmodels.Cursor{
			ID: monasteries[len(monasteries)-1].Id,
		}
	}
	return appmodels.ToPage(monasteries, hasMore, cursor), nil
}

func (s *MonasteryService) CreateMonastery(
	ctx context.Context,
	name string,
) (*models.Monastery, error) {
	monastery := &models.Monastery{
		Name: name,
	}
	err := s.repo.Add(ctx, monastery)
	if err != nil {
		return nil, err
	}
	return monastery, nil
}

func (s *MonasteryService) CheckWriteAccess(ctx context.Context, monasteryID, userID int64) error {
	allowed, err := s.repo.CanWrite(ctx, monasteryID, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return apperr.NewForbidden("monastery.forbidden")
	}
	return nil
}

func (s *MonasteryService) GrantManager(ctx context.Context, monasteryID, userID int64) error {
	return s.repo.GrantManager(ctx, monasteryID, userID)
}

func (s *MonasteryService) RevokeManager(ctx context.Context, monasteryID, userID int64) error {
	return s.repo.RevokeManager(ctx, monasteryID, userID)
}
