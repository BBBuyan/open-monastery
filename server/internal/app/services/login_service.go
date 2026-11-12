package services

import (
	"context"
	"errors"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
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

func (s *LoginService) Login(ctx context.Context, userID int64, password string) (string, error) {
	if userID <= 0 || password == "" {
		return "", ErrInvalidCredentials
	}

	hash, err := s.repo.FindPasswordHash(ctx, userID)
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
