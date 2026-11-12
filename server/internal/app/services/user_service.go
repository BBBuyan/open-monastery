package services

import (
	"context"
	"errors"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
)

type UserService struct {
	repo           abstractions.UserRepository
	passwordHasher abstractions.PasswordHasher
}

func NewUserService(
	repo abstractions.UserRepository,
	passwordHasher abstractions.PasswordHasher,
) *UserService {
	return &UserService{
		repo:           repo,
		passwordHasher: passwordHasher,
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
	user, err := models.NewUser(firstName, lastName, email, phone)
	switch {
	case errors.Is(err, models.ErrContactRequired):
		return nil, apperr.NewInvalid("empty email and phone number", err)
	case errors.Is(err, models.ErrEmailTooLong):
		return nil, apperr.NewInvalidWithMsg(
			"email is too long", "email must be less than 254 characters", err,
		)
	case errors.Is(err, models.ErrEmailNotParseable):
		return nil, apperr.NewInvalid("email not parseable", err)
	case errors.Is(err, models.ErrInvalidEmail):
		return nil, apperr.NewInvalid("invalid email", err)
	case errors.Is(err, models.ErrPhoneNumberTooLong):
		return nil, apperr.NewInvalidWithMsg(
			"phone number is too long", "phone number must be less than 12 characters", err,
		)
	case errors.Is(err, models.ErrInvalidPhoneNumber):
		return nil, apperr.NewInvalidWithMsg(
			"invalid phone number", "phone number is not mongolian", err,
		)
	default:
		return user, err
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
	email *string,
	phone *string,
	password string,
) (*models.User, error) {

	user, err := createBaseUser(firstName, lastName, email, phone)
	if err != nil {
		return nil, err
	}
	err = s.addUser(ctx, user, password)
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
	password string,
) (*models.User, error) {
	user, err := createBaseUser(firstName, lastName, email, phone)
	if err != nil {
		return nil, err
	}

	monk, err := models.NewMonk(user, monasteryId)
	if err != nil {
		return nil, err
	}
	err = s.addUser(ctx, monk, password)
	if err != nil {
		return nil, err
	}

	return monk, nil
}

func (s *UserService) addUser(ctx context.Context, user *models.User, password string) error {
	if password == "" {
		return apperr.NewBadRequest("password.required", nil)
	}
	return s.repo.Add(ctx, user, s.passwordHasher.HashPassword(password))
}
