package dtos

type CategoriesURI struct {
	MonasteryId int64 `uri:"id" binding:"required,min=1"`
}

type CategoryURI struct {
	CategoriesURI
	CategoryId int64 `uri:"categoryId" binding:"required,min=1"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type AddBookToCategoryRequest struct {
	BookId int64 `json:"book_id" binding:"required,gt=0"`
}
