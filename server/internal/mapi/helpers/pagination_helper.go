package helpers

import (
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"github.com/gin-gonic/gin"
)

func BookListOptions(c *gin.Context) (appmodels.ListOptions, error) {
	options := appmodels.ListOptions{Limit: 20, Search: c.Query("q")}
	if value, exists := c.GetQueryArray("limit"); exists {
		limit, err := strconv.Atoi(value[0])
		if err != nil || limit <= 0 {
			return options, apperr.NewBadRequest("invalid limit", err)
		}
		options.Limit = limit
	}
	if value, exists := c.GetQueryArray("after_id"); exists {
		id, err := strconv.ParseInt(value[0], 10, 64)
		if err != nil || id < 0 {
			return options, apperr.NewBadRequest("invalid after_id", err)
		}
		options.AfterId = id
	}
	return options.Normalize()
}
