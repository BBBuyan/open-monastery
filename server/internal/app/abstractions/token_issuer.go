package abstractions

import "example.com/open-monastery/internal/domain/models"

type AccessIdentity struct {
	UserID int64
	Roles  []models.UserType
}

type TokeIssuer interface {
	IssueAccessToken(userID int64, userType models.UserType) (string, error)
	ValidateToken(tokenStr string) (AccessIdentity, error)
}
