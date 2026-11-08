package dtos

import (
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

// BookPageResponse documents the JSON shape of Paginated[*models.Book].
type BookPageResponse struct {
	Data       []*models.Book     `json:"data"`
	Pagination appmodels.PageInfo `json:"pagination"`
}
