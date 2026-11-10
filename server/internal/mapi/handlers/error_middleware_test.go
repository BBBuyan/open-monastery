package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/mapi/dtos"
	"github.com/gin-gonic/gin"
)

func TestErrorHandlerBindingValidation(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.POST("/test", func(c *gin.Context) {
		var req dtos.AddBookToCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(err)
			return
		}
		c.Status(http.StatusNoContent)
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"book_id":0}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), `"title":"validation.error"`) {
		t.Fatalf("expected validation problem with status 400, got %d: %s", response.Code, response.Body.String())
	}
}
