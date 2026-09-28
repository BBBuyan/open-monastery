package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"example.com/open-monastery/internal/apperror"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if appErr, ok := errors.AsType[*apperror.Error](err); ok {
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

func statusFor(code apperror.Code) int {
	switch code {
	case apperror.BadRequest:
		return http.StatusBadRequest
	case apperror.NotFound:
		return http.StatusNotFound
	case apperror.Invalid:
		return http.StatusUnprocessableEntity
	case apperror.Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
