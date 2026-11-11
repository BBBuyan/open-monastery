package services

import (
	"example.com/open-monastery/internal/app/abstractions"
)

type LoginService struct {
	TokenIssuer abstractions.TokeIssuer
}

func NewLoginService(tokenIssuer abstractions.TokeIssuer) *LoginService {
	return &LoginService{
		TokenIssuer: tokenIssuer,
	}
}

func Login(
	userId int64,
	pass string,
) {

}
