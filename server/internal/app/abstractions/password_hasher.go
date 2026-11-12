package abstractions

type PasswordHasher interface {
	HashPassword(password string) string
	VerifyPassword(password string, stored string) bool
}
