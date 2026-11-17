package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/domain/models"
	"example.com/open-monastery/internal/mapi/apimodels"
	"github.com/gin-gonic/gin"
)

const identityKey = "authenticatedIdentity"

func RequireAuth(
	validateToken func(string) (abstractions.AccessIdentity, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.Fields(c.GetHeader("Authorization"))
		if len(header) != 2 || !strings.EqualFold(header[0], "Bearer") {
			abortUnauthorized(c)
			return
		}

		identity, err := validateToken(header[1])
		if err != nil || identity.UserID <= 0 {
			abortUnauthorized(c)
			return
		}

		c.Set(identityKey, identity)
		c.Next()
	}
}

func RequireRole(allowed ...models.UserType) func(*gin.Context) error {
	return func(c *gin.Context) error {
		identity, ok := authenticatedIdentity(c)
		if !ok {
			abortUnauthorized(c)
			return nil
		}

		for _, role := range allowed {
			if slices.Contains(identity.Roles, role) {
				return nil
			}
		}
		return apperr.NewForbidden("user_type.forbidden")
	}
}

// UserID returns the authenticated user ID set by RequireAuth.
func UserID(c *gin.Context) (int64, bool) {
	identity, ok := authenticatedIdentity(c)
	return identity.UserID, ok
}

func authenticatedIdentity(c *gin.Context) (abstractions.AccessIdentity, bool) {
	value, _ := c.Get(identityKey)
	identity, ok := value.(abstractions.AccessIdentity)
	return identity, ok && identity.UserID > 0
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
