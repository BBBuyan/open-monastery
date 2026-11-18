package appmodels

import "example.com/open-monastery/internal/domain/enums"

type Actor struct {
	Id          int64
	MonasteryId *int64
	Role        enums.UserType
}
