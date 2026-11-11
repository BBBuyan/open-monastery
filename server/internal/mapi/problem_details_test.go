package mapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

func TestWithProblemDetailsBindingValidation(t *testing.T) {
	router := gin.New()
	router.POST("/test", withProblemDetails(func(c *gin.Context) error {
		var req dtos.AddBookToCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return apperr.NewBadRequest("invalid.body", err)
		}
		c.Status(http.StatusNoContent)
		return nil
	}))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"book_id":0}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Fatalf("unexpected Content-Type: %s", got)
	}
	body := response.Body.String()
	if response.Code != http.StatusBadRequest || !strings.Contains(body, `"title":"validation.error"`) || !strings.Contains(body, `"errors":[`) {
		t.Fatalf("expected validation problem with status 400, got %d: %s", response.Code, response.Body.String())
	}
}
