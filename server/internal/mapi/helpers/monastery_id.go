package helpers

import (
	"strconv"

	"example.com/open-monastery/internal/app/apperr"
	"github.com/gin-gonic/gin"
)

func MonasteryID(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, apperr.NewBadRequest("invalid monastery id", err)
	}
	return id, nil
}
