package services

import (
	"context"
	"errors"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginService struct {
	repo           abstractions.LoginRepository
	passwordHasher abstractions.PasswordHasher
	tokenIssuer    abstractions.TokeIssuer
	dummyHash      string
}

func NewLoginService(
	repo abstractions.LoginRepository,
	passwordHasher abstractions.PasswordHasher,
	tokenIssuer abstractions.TokeIssuer,
) *LoginService {
	return &LoginService{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenIssuer:    tokenIssuer,
		dummyHash:      passwordHasher.HashPassword(""),
	}
}

func (s *LoginService) Login(ctx context.Context, email string, password string) (string, error) {
	email, err := models.NormalizeEmail(email)
	if err != nil || password == "" {
		return "", ErrInvalidCredentials
	}

	userID, hash, err := s.repo.FindCredentialsByEmail(ctx, email)
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
	return s.tokenIssuer.IssueAccessToken(userID)
}
