package commands

type CreateUserCommand struct {
	FirstName string
	LastName  string
	Email     string
	Phone     *string
	Password  string
}
