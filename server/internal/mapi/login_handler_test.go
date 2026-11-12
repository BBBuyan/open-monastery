package mapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/infra/password"
	"example.com/open-monastery/internal/infra/token"
	"example.com/open-monastery/internal/mapi/apimodels"
	"example.com/open-monastery/internal/mapi/handlers"
	"github.com/gin-gonic/gin"
)

type loginTestRepository struct{ hash string }

func (r loginTestRepository) FindCredentialsByEmail(_ context.Context, email string) (int64, string, error) {
	if email == "user@example.com" {
		return 42, r.hash, nil
	}
	return 0, "", apperr.NewNotFound("credentials not found", nil)
}

func TestLoginByEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hasher := password.NewArgonPasswordHasher()
	issuer := token.NewJwtIssuer(bytes.Repeat([]byte{1}, 32))
	issuer.ValidMinute = 15
	repo := loginTestRepository{hash: hasher.HashPassword("example password")}
	handler := handlers.NewLoginHandler(services.NewLoginService(repo, hasher, issuer))
	router := gin.New()
	router.POST("/api/v1/auth/login", withProblemDetails(handler.Login))

	for _, tc := range []struct {
		name, body string
		status     int
	}{
		{"valid credentials", `{"email":" user@example.com ","password":"example password"}`, http.StatusOK},
		{"wrong password", `{"email":"user@example.com","password":"wrong"}`, http.StatusUnauthorized},
		{"unknown email", `{"email":"missing@example.com","password":"wrong"}`, http.StatusUnauthorized},
		{"missing email", `{"user_id":42,"password":"example password"}`, http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if tc.status != http.StatusOK {
				return
			}
			var result apimodels.LoginResponse
			if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			subject, err := issuer.ValidateToken(result.AccessToken)
			if err != nil || subject == nil || *subject != "42" {
				t.Fatalf("expected a valid JWT with subject 42, got %v", err)
			}
		})
	}
}
