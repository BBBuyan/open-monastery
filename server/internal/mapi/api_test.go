package mapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/domain/models"
	"example.com/open-monastery/internal/infra/token"
	"example.com/open-monastery/internal/mapi/middlewares"
	"github.com/gin-gonic/gin"
)

func TestMonasteryRoutes(t *testing.T) {
	api := RegisterAPI(nil, token.NewJwtIssuer([]byte("test-signing-key")))
	const base = "/api/v1"
	for _, tc := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/monasteries/0/books", http.StatusBadRequest},
		{http.MethodGet, "/monasteries/0/books/2", http.StatusBadRequest},
		{http.MethodGet, "/monasteries/0/books/2/categories", http.StatusBadRequest},
		{http.MethodPost, "/monasteries/1/books", http.StatusUnauthorized},
		{http.MethodGet, "/monasteries/0/categories", http.StatusBadRequest},
		{http.MethodPost, "/monasteries/1/categories", http.StatusUnauthorized},
		{http.MethodGet, "/monasteries/0/categories/2/books", http.StatusBadRequest},
		{http.MethodPost, "/monasteries/1/categories/2/books", http.StatusUnauthorized},
		{http.MethodGet, "/monasteries/1/readings", http.StatusUnauthorized},
		{http.MethodPost, "/monasteries/1/readings", http.StatusUnauthorized},
		{http.MethodPost, "/monasteries/1/readings/2/claim", http.StatusUnauthorized},
		{http.MethodPut, "/monasteries/1/readings/2/assignment", http.StatusUnauthorized},
		{http.MethodPost, "/monasteries", http.StatusUnauthorized},
		{http.MethodPut, "/monasteries/1/managers/2", http.StatusUnauthorized},
		{http.MethodDelete, "/monasteries/1/managers/2", http.StatusUnauthorized},
		{http.MethodGet, "/books", http.StatusNotFound},
		{http.MethodGet, "/categories", http.StatusNotFound},
		{http.MethodGet, "/readings", http.StatusNotFound},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			api.ServeHTTP(response, httptest.NewRequest(tc.method, base+tc.path, nil))
			if response.Code != tc.status {
				t.Fatalf("got status %d, want %d: %s", response.Code, tc.status, response.Body)
			}
		})
	}
}

func TestMonasteryManagerAccess(t *testing.T) {
	issuer := token.NewJwtIssuer([]byte("test-signing-key"))
	issuer.ValidMinute = 15
	api := gin.New()
	public := api.Group("/api/v1")
	authenticated := public.Group("", middlewares.RequireAuth(issuer.ValidateToken))
	monasteryService := services.NewMonasteryService(testMonasteryRepository{})
	addBook(nil, public, authenticated, monasteryService)
	addCategory(nil, public, authenticated, monasteryService)
	addReadings(nil, authenticated, monasteryService)

	for _, tc := range []struct {
		id     string
		userID int64
		role   models.UserType
		status int
	}{
		// An authorized request reaches the handler and fails on its empty body.
		{"7", 2, models.Manager, http.StatusBadRequest},
		{"8", 2, models.Manager, http.StatusForbidden},
		{"7", 3, models.Manager, http.StatusForbidden},
		{"7", 2, models.Monk, http.StatusForbidden},
		{"8", 4, models.Admin, http.StatusBadRequest},
		{"0", 3, models.Admin, http.StatusBadRequest},
	} {
		accessToken, err := issuer.IssueAccessToken(tc.userID, tc.role)
		if err != nil {
			t.Fatal(err)
		}
		for _, route := range []struct{ method, path string }{
			{http.MethodPost, "/books"},
			{http.MethodPost, "/categories"},
			{http.MethodPost, "/categories/5/books"},
			{http.MethodPut, "/readings/5/assignment"},
		} {
			request := httptest.NewRequest(route.method, "/api/v1/monasteries/"+tc.id+route.path, nil)
			request.Header.Set("Authorization", "Bearer "+accessToken)
			response := httptest.NewRecorder()
			api.ServeHTTP(response, request)
			if response.Code != tc.status {
				t.Errorf("%s %s, user %d, role %s: got %d, want %d",
					route.method, request.URL.Path, tc.userID, tc.role, response.Code, tc.status)
			}
		}
	}
}

type testMonasteryRepository struct {
	abstractions.MonasteryRepository
}

func (testMonasteryRepository) CanWrite(_ context.Context, monasteryID, userID int64) (bool, error) {
	return userID == 4 || (userID == 2 && monasteryID == 7), nil
}

func TestOnlyPlatformAdminCanChangeManagerGrants(t *testing.T) {
	issuer := token.NewJwtIssuer([]byte("test-signing-key"))
	issuer.ValidMinute = 15
	api := RegisterAPI(nil, issuer)
	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		for _, role := range []models.UserType{models.Monk, models.Admin} {
			accessToken, err := issuer.IssueAccessToken(2, role)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(method, "/api/v1/monasteries/7/managers/0", nil)
			request.Header.Set("Authorization", "Bearer "+accessToken)
			response := httptest.NewRecorder()
			api.ServeHTTP(response, request)
			want := http.StatusForbidden
			if role == models.Admin {
				want = http.StatusBadRequest
			}
			if response.Code != want {
				t.Errorf("%s, role %s: got %d, want %d", method, role, response.Code, want)
			}
		}
	}
}
