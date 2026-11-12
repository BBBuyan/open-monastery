package mapi

import (
	"errors"
	"net/http"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/mapi/apimodels"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const internalServerErrorTitle = "internal server error"
const validationErrorTitle = "validation.error"

func withProblemDetails(handle func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handle(c)
		if err == nil || c.Writer.Written() {
			return
		}

		status, body := problemFromError(err, c.Request.URL.Path)
		c.Header("Content-Type", "application/problem+json")
		c.AbortWithStatusJSON(status, body)
	}
}

func problemFromError(err error, path string) (int, any) {
	if vErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
		var problem apimodels.ValidationProblemDetails
		problem.Status = http.StatusBadRequest
		problem.Title = new(validationErrorTitle)
		problem.Instance = &path
		problem.Errors = make([]string, 0, len(vErr))
		for _, e := range vErr {
			problem.Errors = append(problem.Errors, e.Error())
		}
		return problem.Status, problem
	}

	problem := apimodels.ProblemDetails{Instance: &path}
	if aErr, ok := errors.AsType[*apperr.AppErr](err); ok {
		problem.Status = statusFor(aErr.Type)
		problem.Title = &aErr.Code
		problem.Detail = aErr.Message
		return problem.Status, problem
	}

	problem.Status = http.StatusInternalServerError
	problem.Title = new(internalServerErrorTitle)
	return problem.Status, problem
}

func statusFor(errType apperr.AppErrType) int {
	switch errType {
	case apperr.BadRequest:
		return http.StatusBadRequest
	case apperr.Unauthorized:
		return http.StatusUnauthorized
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
