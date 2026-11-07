package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if appErr, ok := errors.AsType[*apperr.AppErr](err); ok {
			c.AbortWithStatusJSON(statusFor(appErr.Code), gin.H{
				"code":  appErr.Code,
				"error": appErr.Message,
			})
			return
		}

		slog.ErrorContext(c.Request.Context(), "request failed", "error", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":  "internal_error",
			"error": "internal server error",
		})
	}
}

func statusFor(code apperr.AppErrCode) int {
	switch code {
	case apperr.BadRequest:
		return http.StatusBadRequest
	case apperr.NotFound:
		return http.StatusNotFound
	case apperr.Invalid:
		return http.StatusUnprocessableEntity
	case apperr.Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
