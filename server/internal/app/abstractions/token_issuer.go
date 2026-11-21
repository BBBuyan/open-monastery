package abstractions

import (
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
)

type TokenIssuer interface {
	IssueAccessToken(userID int64, userType enums.UserType, monasteryId *int64) (string, error)
	ValidateToken(tokenStr string) (appmodels.Actor, error)
}
