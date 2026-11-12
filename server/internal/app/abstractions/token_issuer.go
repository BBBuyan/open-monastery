package abstractions

type TokeIssuer interface {
	IssueAccessToken(userId int64) (string, error)
	ValidateToken(tokenStr string) (*string, error)
}
