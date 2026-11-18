package commands

type CreateBookCommand struct {
	Name        string
	Desc        string
	Price       uint32
	UserId      int64
	MonasteryID int64
}
