package mapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/domain/models"
	"example.com/open-monastery/internal/mapi/handlers"
	"github.com/gin-gonic/gin"
)

type monkTestRepository struct {
	saved *models.User
	err   error
}

func (r *monkTestRepository) FindById(context.Context, int64) (*models.User, error) {
	return r.saved, nil
}
func (r *monkTestRepository) Add(_ context.Context, user *models.User) error {
	if r.err != nil {
		return r.err
	}
	user.Id = 42
	r.saved = user
	return nil
}

func TestRegisterMonk(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, body string
		repoErr    error
		status     int
	}{
		{name: "email", body: `{"first_name":"A","last_name":"B","email":" a@example.com ","monastery_id":7}`, status: 201},
		{name: "phone", body: `{"first_name":"A","last_name":"B","phone_number":"+97699112233","monastery_id":7}`, status: 201},
		{name: "missing monastery", body: `{"first_name":"A","last_name":"B","email":"a@example.com"}`, status: 400},
		{name: "negative monastery", body: `{"first_name":"A","last_name":"B","email":"a@example.com","monastery_id":-1}`, status: 400},
		{name: "missing name", body: `{"email":"a@example.com","monastery_id":7}`, status: 400},
		{name: "missing contact", body: `{"first_name":"A","last_name":"B","monastery_id":7}`, status: 400},
		{name: "invalid contact", body: `{"first_name":"A","last_name":"B","email":"invalid","monastery_id":7}`, status: 422},
		{name: "malformed JSON", body: `{`, status: 400},
		{name: "write failure", body: `{"first_name":"A","last_name":"B","email":"a@example.com","monastery_id":7}`, repoErr: errors.New("write failed"), status: 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &monkTestRepository{err: tc.repoErr}
			router := gin.New()

			router.POST("/monks", withProblemDetails(handlers.NewUserHandler(services.NewUserService(repo)).RegisterMonk))
			request := httptest.NewRequest("POST", "/monks", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if tc.status != 201 {
				if repo.saved != nil {
					t.Fatal("invalid request persisted a user")
				}
				return
			}
			var user models.User
			if err := json.Unmarshal(response.Body.Bytes(), &user); err != nil {
				t.Fatal(err)
			}
			if user.Id != 42 || user.UserType != models.Monk || user.MonasteryId == nil || *user.MonasteryId != 7 {
				t.Fatalf("unexpected monk: %+v", user)
			}
			if repo.saved == nil || repo.saved.UserType != models.Monk || repo.saved.MonasteryId == nil || *repo.saved.MonasteryId != 7 {
				t.Fatal("monk role and monastery were not persisted")
			}
			if tc.name == "email" && (user.Email == nil || *user.Email != "a@example.com" || repo.saved.Email == nil || *repo.saved.Email != "a@example.com") {
				t.Fatal("normalized email was not preserved")
			}
			if tc.name == "phone" && (user.PhoneNumber == nil || *user.PhoneNumber != "+97699112233" || repo.saved.PhoneNumber == nil || *repo.saved.PhoneNumber != "+97699112233") {
				t.Fatal("phone was not preserved")
			}
		})
	}
}

func TestRegisterUserCreatesCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, body := range []string{
		`{"first_name":"A","last_name":"B","email":"a@example.com"}`,
		`{"first_name":"A","last_name":"B","email":"a@example.com","user_type":"monk","monastery_id":7}`,
	} {
		repo := &monkTestRepository{}
		router := gin.New()

		router.POST("/users", withProblemDetails(handlers.NewUserHandler(services.NewUserService(repo)).RegisterUser))
		request := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != 201 {
			t.Fatalf("status = %d: %s", response.Code, response.Body.String())
		}
		var user models.User
		if err := json.Unmarshal(response.Body.Bytes(), &user); err != nil {
			t.Fatal(err)
		}
		if user.UserType != models.Customer || user.MonasteryId != nil {
			t.Fatalf("unexpected user: %+v", user)
		}
		if repo.saved == nil || repo.saved.UserType != models.Customer || repo.saved.MonasteryId != nil {
			t.Fatal("expected customer to be persisted without monastery")
		}
	}
}

func TestRegisterUserRequiresNames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, body := range []string{
		`{"last_name":"B","email":"a@example.com"}`,
		`{"first_name":"A","email":"a@example.com"}`,
		`{"first_name":"","last_name":"B","email":"a@example.com"}`,
		`{"first_name":"A","last_name":null,"email":"a@example.com"}`,
	} {
		repo := &monkTestRepository{}
		router := gin.New()

		router.POST("/users", withProblemDetails(handlers.NewUserHandler(services.NewUserService(repo)).RegisterUser))
		request := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != 400 {
			t.Fatalf("status = %d: %s", response.Code, response.Body.String())
		}
		if repo.saved != nil {
			t.Fatal("request with missing name persisted a user")
		}
	}
}
