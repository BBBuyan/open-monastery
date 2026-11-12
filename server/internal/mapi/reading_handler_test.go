package mapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/domain/models"
	"example.com/open-monastery/internal/infra/token"
	"example.com/open-monastery/internal/mapi/handlers"
	"example.com/open-monastery/internal/mapi/middlewares"
	"github.com/gin-gonic/gin"
)

type readingAuthRepository struct {
	abstractions.ReadingRepository
	saved *models.Reading
}

func (r *readingAuthRepository) FindByIdForCustomer(_ context.Context, id, customerID int64) (*models.Reading, error) {
	if r.saved == nil || r.saved.Id != id || r.saved.CustomerId != customerID {
		return nil, apperr.NewNotFound("reading not found", nil)
	}
	return r.saved, nil
}

func (r *readingAuthRepository) Add(_ context.Context, reading *models.Reading) error {
	reading.Id = 1
	r.saved = reading
	return nil
}

func (*readingAuthRepository) UserExists(context.Context, int64) (bool, error) { return true, nil }
func (*readingAuthRepository) BookExists(context.Context, int64) (bool, error) { return true, nil }

func TestReadingRequesterAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := token.NewJwtIssuer(bytes.Repeat([]byte{1}, 32))
	issuer.ValidMinute = 15
	issue := func(userID int64) string {
		t.Helper()
		accessToken, err := issuer.IssueAccessToken(userID)
		if err != nil {
			t.Fatal(err)
		}
		return "Bearer " + accessToken
	}
	owner, other := issue(42), issue(99)

	repo := &readingAuthRepository{}
	handler := handlers.NewReadingHandler(services.NewReadingService(repo))
	router := gin.New()
	readings := router.Group("/api/v1/readings", middlewares.RequireAuth(issuer.ValidateToken))
	readings.GET("/:id", withProblemDetails(handler.GetReading))
	readings.POST("", withProblemDetails(handler.AddReading))

	for _, tc := range []struct {
		name, method, path, auth, body string
		status                         int
	}{
		{"anonymous create", "POST", "", "", `{"book_id":3}`, http.StatusUnauthorized},
		{"anonymous read", "GET", "/1", "", "", http.StatusUnauthorized},
		{"invalid token", "POST", "", "Bearer invalid", `{"book_id":3}`, http.StatusUnauthorized},
		{"create as requester", "POST", "", owner, `{"book_id":3,"customer_id":99}`, http.StatusCreated},
		{"owner reads", "GET", "/1", owner, "", http.StatusOK},
		{"other user cannot read", "GET", "/1", other, "", http.StatusNotFound},
		{"missing reading", "GET", "/2", owner, "", http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/v1/readings"+tc.path, strings.NewReader(tc.body))
			req.Header.Set("Authorization", tc.auth)
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if tc.status == http.StatusCreated && (repo.saved == nil || repo.saved.CustomerId != 42 || repo.saved.MonkId != nil) {
				t.Fatalf("expected an unassigned reading owned by the token subject, got %+v", repo.saved)
			}
		})
	}
}
