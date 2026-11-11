package abstractions

type TokeIssuer interface {
	IssueAccessToken(userId int64) (string, error)
}
