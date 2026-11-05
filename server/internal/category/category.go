package category

type Category struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryBook struct {
	Id         int64
	CategoryId int64
	BookId     int64
}
