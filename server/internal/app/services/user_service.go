package services

import (
	"context"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/helpers"
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
func createBaseUser(
	firstName string,
	lastName string,
	email *string,
	phone *string,
) (*models.User, error) {
	if email == nil && phone == nil {
		return nil, apperr.NewBadRequest(
			"empty email, phone number",
			"email and phone number is both empty",
			nil,
		)
	}
	if email != nil {
		e, err := helpers.NewEmail(*email)
		if err != nil {
			return nil, err
		}
		email = &e
	}
	if phone != nil {
		p, err := helpers.NewPhoneNumber(*phone)
		if err != nil {
			return nil, err
		}
		phone = &p
	}

	user, err := models.NewUser(
		firstName, lastName,
		email, phone,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
	email *string,
	phone *string,
) (*models.User, error) {

	user, err := createBaseUser(firstName, lastName, email, phone)
	if err != nil {
		return nil, err
	}
	err = s.repo.Add(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateMonk(
	ctx context.Context,
	firstName string,
	lastName string,
	email *string,
	phone *string,
	monasteryId int64,
) (*models.User, error) {
	user, err := createBaseUser(firstName, lastName, email, phone)
	if err != nil {
		return nil, err
	}

	monk, err := models.NewMonk(user, monasteryId)
	if err != nil {
		return nil, err
	}
	err = s.repo.Add(ctx, monk)
	if err != nil {
		return nil, err
	}

	return monk, nil
}
