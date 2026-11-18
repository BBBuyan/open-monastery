package abstractions

import (
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
)

type TokeIssuer interface {
	IssueAccessToken(userID int64, userType enums.UserType) (string, error)
	ValidateToken(tokenStr string) (appmodels.Actor, error)
}
