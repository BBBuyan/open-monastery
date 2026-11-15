package services

import (
	"context"
	"errors"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	repo           abstractions.AuthRepository
	passwordHasher abstractions.PasswordHasher
	tokenIssuer    abstractions.TokeIssuer
	dummyHash      string
}

func NewAuthService(
	repo abstractions.AuthRepository,
	passwordHasher abstractions.PasswordHasher,
	tokenIssuer abstractions.TokeIssuer,
) *AuthService {
	return &AuthService{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenIssuer:    tokenIssuer,
		dummyHash:      passwordHasher.HashPassword(""),
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	email, err := models.NormalizeEmail(email)
	if err != nil || password == "" {
		return "", ErrInvalidCredentials
	}

	userID, hash, userType, err := s.repo.FindCredentialsByEmail(ctx, email)
	if err != nil {
		appErr, ok := errors.AsType[*apperr.AppErr](err)
		if !ok || appErr.Type != apperr.NotFound {
			return "", err
		}
		// Match the password verification cost for an account without credentials.
		s.passwordHasher.VerifyPassword(password, s.dummyHash)
		return "", ErrInvalidCredentials
	}

	if !s.passwordHasher.VerifyPassword(password, hash) {
		return "", ErrInvalidCredentials
	}
	return s.tokenIssuer.IssueAccessToken(userID, userType)
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

func (s *AuthService) CreateUser(
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

func (s *AuthService) CreateMonk(
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

func (s *AuthService) addUser(ctx context.Context, user *models.User, password string) error {
	if password == "" {
		return apperr.NewBadRequest("password.required", nil)
	}
	return s.repo.AddUser(ctx, user, s.passwordHasher.HashPassword(password))
}
