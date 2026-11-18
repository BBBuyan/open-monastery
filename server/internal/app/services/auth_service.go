package services

import (
	"context"
	"errors"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/app/commands"
	"example.com/open-monastery/internal/app/security"
	"example.com/open-monastery/internal/domain/enums"
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

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	email, err := models.NormalizeEmail(email)
	if err != nil || password == "" {
		return "", apperr.NewInvalid("invalid credentials", err)
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

func (s *AuthService) CreateUser(
	ctx context.Context,
	cmd commands.CreateUserCommand,
) (*models.User, error) {
	if err := security.ValidatePassword(cmd.Password); err != nil {
		return nil, err
	}
	passHash := s.passwordHasher.HashPassword(cmd.Password)
	email, err := models.NewEmail(cmd.Email)
	if err != nil {
		return nil, apperr.NewInvalid(err.Error(), err)
	}
	var phone *models.PhoneNumber
	if cmd.Phone != nil {
		phone, err = models.NewPhoneNumber(*cmd.Phone)
		if err != nil {
			return nil, apperr.NewInvalid(err.Error(), err)
		}
	}
	user, err := models.NewUser(
		cmd.FirstName, cmd.LastName,
		*email, phone,
		passHash,
	)
	if err != nil {
		return nil, err
	}
	err = s.repo.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) PromoteToMonk(
	ctx context.Context,
	userId int64,
	actor appmodels.Actor,
) error {
	if actor.Role != enums.Manager {
		return apperr.NewForbidden("forbidden")
	}
	if actor.MonasteryId == nil {
		return apperr.NewBadRequest("manager monastery id is empty", nil)
	}
	s.repo.PromoteToMonk(ctx, userId, *actor.MonasteryId)

	return nil
}

func (s *AuthService) PromoteToManager(
	ctx context.Context,
	monkId int64,
	monasteryId int64,
	actor appmodels.Actor,
) error {
	if actor.Role != enums.Admin {
		return apperr.NewForbidden("forbidden")
	}
	s.repo.PromoteToManager(ctx, monkId, monasteryId)

	return nil
}
