package abstractions

import (
	"example.com/open-monastery/internal/app/appmodels"
)

type TokenIssuer interface {
	IssueAccessToken(creds *appmodels.UserCreds) (string, error)
	ValidateToken(tokenStr string) (appmodels.Actor, error)
}
