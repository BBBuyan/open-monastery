package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
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
