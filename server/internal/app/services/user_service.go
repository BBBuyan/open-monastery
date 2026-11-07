package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/domain/models"
)

type UserService struct {
	repo abstractions.UserRepository
}

func NewUserService(repo abstractions.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUser(
	ctx context.Context,
	id int64,
) (*models.User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *UserService) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
) (*models.User, error) {
	user := &models.User{
		FirstName: firstName,
		LastName:  lastName,
		UserType:  models.Customer,
	}

	err := s.repo.Add(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateMonk(
	ctx context.Context,
	firstName string,
	lastName string,
	monasteryId int64,
) (*models.User, error) {
	monk := &models.User{
		FirstName:   firstName,
		LastName:    lastName,
		UserType:    models.Monk,
		MonasteryId: &monasteryId,
	}

	err := s.repo.Add(ctx, monk)

	if err != nil {
		return nil, err
	}

	return monk, nil
}
