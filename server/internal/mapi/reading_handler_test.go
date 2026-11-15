package mapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
	"example.com/open-monastery/internal/infra/token"
	"example.com/open-monastery/internal/mapi/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type readingAuthRepository struct {
	readings []*models.Reading
}

func (r *readingAuthRepository) GetReadingsByUserId(_ context.Context, userID int64, opts appmodels.ListOptions) ([]*models.Reading, error) {
	readings := make([]*models.Reading, 0)
	for _, reading := range r.readings {
		if reading.CustomerId == userID && reading.Id > opts.AfterId {
			readings = append(readings, reading)
			if len(readings) == opts.Limit+1 {
				break
			}
		}
	}
	return readings, nil
}

func (r *readingAuthRepository) Add(_ context.Context, reading *models.Reading) error {
	reading.Id = int64(len(r.readings) + 1)
	r.readings = append(r.readings, reading)
	return nil
}

func (r *readingAuthRepository) Claim(_ context.Context, readingID, monkID int64) error {
	for _, reading := range r.readings {
		if reading.Id != readingID {
			continue
		}
		if reading.MonkId != nil {
			return apperr.NewConflict("reading.already_claimed", nil)
		}
		reading.AssignToMonk(monkID)
		return nil
	}
	return apperr.NewNotFound("reading.not_found", nil)
}

func (*readingAuthRepository) UserExists(context.Context, int64) (bool, error) { return true, nil }
func (*readingAuthRepository) BookExists(context.Context, int64) (bool, error) { return true, nil }

func TestReadingRequesterAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := token.NewJwtIssuer(bytes.Repeat([]byte{1}, 32))
	issuer.ValidMinute = 15
	issue := func(userID int64) string {
		t.Helper()
		accessToken, err := issuer.IssueAccessToken(userID, models.Customer)
		if err != nil {
			t.Fatal(err)
		}
		return "Bearer " + accessToken
	}
	owner, other := issue(42), issue(99)

	repo := &readingAuthRepository{readings: []*models.Reading{
		{Id: 1, BookId: 3, CustomerId: 42},
		{Id: 2, BookId: 3, CustomerId: 99},
	}}
	router := gin.New()
	registerReadings(router.Group("/api/v1"), repo, middlewares.RequireAuth(issuer.ValidateToken))

	cursor, err := appmodels.EncodeCursor(appmodels.Cursor{ID: 1})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name, method, path, auth, body string
		status                         int
		ids                            []int64
		nextCursor                     string
	}{
		{"anonymous create", "POST", "", "", `{"book_id":3}`, http.StatusUnauthorized, nil, ""},
		{"anonymous list", "GET", "", "", "", http.StatusUnauthorized, nil, ""},
		{"invalid token", "POST", "", "Bearer invalid", `{"book_id":3}`, http.StatusUnauthorized, nil, ""},
		{"create as requester", "POST", "", owner, `{"book_id":3,"customer_id":99}`, http.StatusCreated, nil, ""},
		{"owner lists all readings", "GET", "", owner, "", http.StatusOK, []int64{1, 3}, ""},
		{"first page", "GET", "?limit=1", owner, "", http.StatusOK, []int64{1}, cursor},
		{"last page", "GET", "?limit=1&cursor=" + cursor, owner, "", http.StatusOK, []int64{3}, ""},
		{"list stays scoped to token subject", "GET", "?user_id=42", other, "", http.StatusOK, []int64{2}, ""},
		{"empty list", "GET", "", issue(100), "", http.StatusOK, []int64{}, ""},
		{"invalid limit", "GET", "?limit=0", owner, "", http.StatusBadRequest, nil, ""},
		{"invalid cursor", "GET", "?cursor=!", owner, "", http.StatusBadRequest, nil, ""},
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
			if tc.status == http.StatusCreated {
				saved := repo.readings[len(repo.readings)-1]
				if saved.CustomerId != 42 || saved.MonkId != nil {
					t.Fatalf("expected an unassigned reading owned by the token subject, got %+v", saved)
				}
			}
			if tc.status == http.StatusOK {
				var page appmodels.Paginated[models.Reading]
				if err := json.Unmarshal(response.Body.Bytes(), &page); err != nil {
					t.Fatal(err)
				}
				ids := make([]int64, 0, len(page.Data))
				for _, reading := range page.Data {
					ids = append(ids, reading.Id)
				}
				if page.Data == nil || !slices.Equal(ids, tc.ids) {
					t.Fatalf("got data %v, want reading IDs %v", page.Data, tc.ids)
				}
				if page.Pagination.HasMore != (tc.nextCursor != "") ||
					(tc.nextCursor == "" && page.Pagination.NextCursor != nil) ||
					(tc.nextCursor != "" && (page.Pagination.NextCursor == nil || *page.Pagination.NextCursor != tc.nextCursor)) {
					t.Fatalf("expected next cursor %q and matching hasMore, got %s", tc.nextCursor, response.Body.String())
				}
			}
		})
	}
}

func TestClaimReadingRequiresMonk(t *testing.T) {
	gin.SetMode(gin.TestMode)
	signingKey := bytes.Repeat([]byte{1}, 32)
	issuer := token.NewJwtIssuer(signingKey)
	issuer.ValidMinute = 15
	reading := &models.Reading{Id: 1, CustomerId: 10}
	repo := &readingAuthRepository{readings: []*models.Reading{reading}}
	router := gin.New()
	registerReadings(router.Group("/api/v1"), repo, middlewares.RequireAuth(issuer.ValidateToken))

	for _, tc := range []struct {
		name      string
		readingID string
		userID    int64
		role      models.UserType
		status    int
		assignee  int64
	}{
		{"anonymous", "1", 0, "", http.StatusUnauthorized, 0},
		{"legacy token without roles", "1", 20, "", http.StatusForbidden, 0},
		{"customer", "1", 10, models.Customer, http.StatusForbidden, 0},
		{"admin", "1", 40, models.Admin, http.StatusForbidden, 0},
		{"invalid reading ID", "0", 20, models.Monk, http.StatusBadRequest, 0},
		{"missing reading", "2", 20, models.Monk, http.StatusNotFound, 0},
		{"monk claims reading", "1", 20, models.Monk, http.StatusNoContent, 20},
		{"another monk cannot take over", "1", 30, models.Monk, http.StatusConflict, 20},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/readings/"+tc.readingID+"/claim",
				strings.NewReader(`{"monk_id":30}`))
			if tc.userID != 0 {
				accessToken, err := issuer.IssueAccessToken(tc.userID, tc.role)
				if err != nil {
					t.Fatal(err)
				}
				if tc.role == "" {
					parsed, _, err := jwt.NewParser().ParseUnverified(accessToken, jwt.MapClaims{})
					if err != nil {
						t.Fatal(err)
					}
					delete(parsed.Claims.(jwt.MapClaims), "roles")
					accessToken, err = parsed.SignedString(signingKey)
					if err != nil {
						t.Fatal(err)
					}
				}
				req.Header.Set("Authorization", "Bearer "+accessToken)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if tc.assignee == 0 && reading.MonkId != nil ||
				tc.assignee != 0 && (reading.MonkId == nil || *reading.MonkId != tc.assignee) {
				t.Fatalf("expected assignee %d, got %+v", tc.assignee, reading)
			}
		})
	}
}
