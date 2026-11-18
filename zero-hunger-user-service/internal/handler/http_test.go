package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/user-service/internal/domain"
	"github.com/zero-hunger/user-service/internal/usecase"
)

type profileRepository struct {
	domain.UserRepository
	calls int
}

func (r *profileRepository) GetByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	r.calls++
	return domain.User{ID: id, Email: "private@example.test", Phone: "private-phone"}, nil
}

func TestGetUserAuthorization(t *testing.T) {
	for _, tc := range []struct {
		name, role      string
		self, anonymous bool
		want            int
	}{
		{name: "other recipient", role: "recipient", want: http.StatusForbidden},
		{name: "other donor", role: "donor", want: http.StatusForbidden},
		{name: "self", role: "recipient", self: true, want: http.StatusOK},
		{name: "admin", role: "admin", want: http.StatusOK},
		{name: "anonymous", anonymous: true, want: http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			caller, target := uuid.New(), uuid.New()
			if tc.self {
				target = caller
			}
			repo := &profileRepository{}
			svc := usecase.NewUserService(repo, nil, "test-only-secret", time.Hour, time.Hour)
			h := NewHTTPHandler(svc)
			e := echo.New()
			e.GET("/api/v1/users/:user_id", h.GetUser)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+target.String(), nil)
			if !tc.anonymous {
				raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": caller.String(), "role": tc.role, "exp": time.Now().Add(time.Hour).Unix(),
				}).SignedString([]byte("test-only-secret"))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Authorization", "Bearer "+raw)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
			}
			if tc.want != http.StatusOK && (repo.calls != 0 || strings.Contains(rec.Body.String(), "private")) {
				t.Fatal("denied access must not read or expose another user's contact details")
			}
		})
	}
}
