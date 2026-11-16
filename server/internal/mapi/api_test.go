package mapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/open-monastery/internal/infra/token"
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
		{http.MethodPost, "/monasteries/0/books", http.StatusBadRequest},
		{http.MethodGet, "/monasteries/0/categories", http.StatusBadRequest},
		{http.MethodPost, "/monasteries/0/categories", http.StatusBadRequest},
		{http.MethodGet, "/monasteries/0/categories/2/books", http.StatusBadRequest},
		{http.MethodPost, "/monasteries/0/categories/2/books", http.StatusBadRequest},
		{http.MethodGet, "/monasteries/1/readings", http.StatusUnauthorized},
		{http.MethodPost, "/monasteries/1/readings", http.StatusUnauthorized},
		{http.MethodPost, "/monasteries/1/readings/2/claim", http.StatusUnauthorized},
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
