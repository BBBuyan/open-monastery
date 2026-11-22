package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"example.com/open-monastery/internal/api/apimodels"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
	"github.com/gin-gonic/gin"
)

const actorKey = "authenticatedActor"

func RequireAuth(
	validateToken func(string) (appmodels.Actor, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.Fields(c.GetHeader("Authorization"))
		if len(header) != 2 || !strings.EqualFold(header[0], "Bearer") {
			abortUnauthorized(c)
			return
		}

		actor, err := validateToken(header[1])
		if err != nil || actor.Id <= 0 {
			abortUnauthorized(c)
			return
		}

		c.Set(actorKey, actor)
		c.Next()
	}
}

func RequireRole(allowed ...enums.UserType) func(*gin.Context) error {
	return func(c *gin.Context) error {
		actor, err := AuthenticatedActor(c)
		if err != nil {
			abortUnauthorized(c)
			return nil
		}

		if slices.Contains(allowed, actor.Role) {
			return nil
		}
		return apperr.NewForbidden("user_type.forbidden")
	}
}

func AuthenticatedActor(c *gin.Context) (*appmodels.Actor, error) {
	value, _ := c.Get(actorKey)
	actor, ok := value.(appmodels.Actor)
	if !ok || actor.Id < 0 {
		return &actor, apperr.NewInvalid("actor unparseable", nil)
	}
	return &actor, nil
}

func abortUnauthorized(c *gin.Context) {
	title := "authentication.required"
	path := c.Request.URL.Path

	c.Header("WWW-Authenticate", "Bearer")
	c.Header("Content-Type", "application/problem+json")
	c.AbortWithStatusJSON(http.StatusUnauthorized, apimodels.ProblemDetails{
		Status:   http.StatusUnauthorized,
		Title:    &title,
		Instance: &path,
	})
}
