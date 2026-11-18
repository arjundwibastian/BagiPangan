package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/food-service/internal/domain"
)

type nearbyUseCase struct {
	domain.FoodListingUseCase
	err    error
	caller uuid.UUID
}

func (u *nearbyUseCase) SearchNearbyFoodListing(_ context.Context, userID, _ uuid.UUID) (*[]domain.FoodListing, error) {
	u.caller = userID
	listings := make([]domain.FoodListing, 0)
	return &listings, u.err
}

func TestNearbyHTTPErrorsAndEmptyResults(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"empty search is successful", nil, http.StatusOK},
		{"another user's request", domain.ErrForbidden, http.StatusForbidden},
		{"missing request", domain.ErrNotFound, http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			uc := &nearbyUseCase{err: tc.err}
			h := NewFoodListingHandler(uc)
			e := echo.New()
			caller := uuid.New()
			e.GET("/nearby", h.GetNearbyFoodListingHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error { c.Set("user_id", caller.String()); return next(c) }
			})
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nearby?request_id="+uuid.NewString(), nil))
			if rec.Code != tc.status || uc.caller != caller {
				t.Fatalf("HTTP %d, caller %s: %s", rec.Code, uc.caller, rec.Body.String())
			}
			if tc.err == nil {
				var body struct {
					Data []domain.FoodListing `json:"responseData"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Data == nil || len(body.Data) != 0 {
					t.Fatalf("empty result must be [], got %s", rec.Body.String())
				}
			}
		})
	}
}
