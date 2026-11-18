package validator

import (
	"testing"
	"time"

	"github.com/zero-hunger/food-service/internal/dto"
)

func TestZeroCoordinatesAreValid(t *testing.T) {
	zero := 0.0
	listing := dto.CreateFoodListingRequest{Title: "Test food", Latitude: &zero, Longitude: &zero, FoodType: "cooked_food", Quantity: 1,
		AvailableFrom: time.Now(), AvailableUntil: time.Now().Add(time.Hour)}
	if err := NewValidator().Validate(listing); err != nil {
		t.Fatalf("zero latitude/longitude is a real location: %v", err)
	}
	listing.Latitude = nil
	if err := NewValidator().Validate(listing); err == nil {
		t.Fatal("missing latitude must still be rejected")
	}
	listing.Latitude, listing.Longitude = &zero, nil
	if err := NewValidator().Validate(listing); err == nil {
		t.Fatal("missing longitude must still be rejected")
	}
}
