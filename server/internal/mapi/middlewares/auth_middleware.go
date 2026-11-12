package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/open-monastery/internal/mapi/apimodels"
	"github.com/gin-gonic/gin"
)

const userIDKey = "authenticatedUserID"

func RequireAuth(validateToken func(string) (*string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.Fields(c.GetHeader("Authorization"))
		if len(header) != 2 || !strings.EqualFold(header[0], "Bearer") {
			abortUnauthorized(c)
			return
		}

		subject, err := validateToken(header[1])
		if err != nil || subject == nil {
			abortUnauthorized(c)
			return
		}

		userID, err := strconv.ParseInt(*subject, 10, 64)
		if err != nil || userID <= 0 {
			abortUnauthorized(c)
			return
		}

		c.Set(userIDKey, userID)
		c.Next()
	}
}

// UserID returns the authenticated user ID set by RequireAuth.
func UserID(c *gin.Context) (int64, bool) {
	value, _ := c.Get(userIDKey)
	userID, ok := value.(int64)
	return userID, ok && userID > 0
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
