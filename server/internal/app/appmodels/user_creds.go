package appmodels

import "example.com/open-monastery/internal/domain/enums"

type UserCreds struct {
	UserId       int64
	MonasteryId  *int64
	UserType     enums.UserType
	PasswordHash string
	Email        string
}
