package dtos

type CreateMonasteryRequest struct {
	Name string `json:"name" binding:"required"`
}

type MonasteriesURI struct {
	MonasteryId int64 `uri:"id" binding:"required,min=1"`
}
