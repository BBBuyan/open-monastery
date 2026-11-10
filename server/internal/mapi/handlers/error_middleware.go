package handlers

import (
	"errors"
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/mapi/apimodels"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TODO: refine
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		if vErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			var vProblem apimodels.ValidationProblemDetails
			vProblem.Status = http.StatusBadRequest
			vProblem.Title = "validation.error"
			vProblem.Instance = c.Request.URL.Path
			errs := make(map[string]string)
			for _, e := range vErr {
				errs[e.Namespace()] = e.Tag()
			}
			vProblem.Errors = errs

			c.AbortWithStatusJSON(vProblem.Status, vProblem)
			return
		}

		var problem apimodels.ProblemDetails
		if aErr, ok := errors.AsType[*apperr.AppErr](err); ok {
			problem.Status = statusFor(aErr.Type)
			problem.Title = aErr.Code
			problem.Detail = aErr.Message
			problem.Instance = c.Request.URL.Path
			c.AbortWithStatusJSON(problem.Status, problem)
			return
		}

		problem.Status = http.StatusInternalServerError
		problem.Title = "internal.error"
		problem.Detail = "internal server error"
		problem.Instance = c.Request.URL.Path

		c.AbortWithStatusJSON(problem.Status, problem)
	}
}

func statusFor(errType apperr.AppErrType) int {
	switch errType {
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
