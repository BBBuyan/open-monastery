package helpers

import (
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"github.com/gin-gonic/gin"
)

func ToListOpts(c *gin.Context) (appmodels.ListOptions, error) {
	options := appmodels.ListOptions{Limit: 20, Search: c.Query("q")}
	if value, exists := c.GetQueryArray("limit"); exists {
		limit, err := strconv.Atoi(value[0])
		if err != nil || limit <= 0 {
			return options, apperr.NewBadRequestWithMsg(
				"invalid limit",
				"limit must be greater than 0",
				err,
			)
		}
		options.Limit = limit
	}
	if value, exists := c.GetQueryArray("cursor"); exists {
		decodedCursor, err := appmodels.DecodeCursor(value[0])
		if err != nil {
			return options, apperr.NewBadRequestWithMsg(
				"invalid cursor",
				"cursor must be not null",
				err,
			)
		}
		options.AfterId = decodedCursor.ID
	}
	return options.Normalize()
}
