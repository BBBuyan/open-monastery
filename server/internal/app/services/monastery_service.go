package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
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
	userId int64,
	actor appmodels.Actor,
) (*models.Monastery, error) {
	if actor.Role != enums.Admin {
		return nil, apperr.NewForbidden("forbidden")
	}
	monastery := models.NewMonastery(name, userId)
	err := s.repo.Add(ctx, monastery)
	if err != nil {
		return nil, err
	}
	return monastery, nil
}
