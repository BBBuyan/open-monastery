package handlers

import (
	"strconv"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/book"
	"github.com/gin-gonic/gin"
)

func bookListOptions(c *gin.Context) (book.ListOptions, error) {
	options := book.ListOptions{Limit: 20, Search: c.Query("q")}
	if value, exists := c.GetQueryArray("limit"); exists {
		limit, err := strconv.Atoi(value[0])
		if err != nil || limit <= 0 {
			return options, apperror.New(apperror.BadRequest, "invalid limit", err)
		}
		options.Limit = limit
	}
	if value, exists := c.GetQueryArray("after_id"); exists {
		id, err := strconv.ParseInt(value[0], 10, 64)
		if err != nil || id < 0 {
			return options, apperror.New(apperror.BadRequest, "invalid after_id", err)
		}
		options.AfterId = id
	}
	return options.Normalize()
}
