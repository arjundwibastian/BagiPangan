package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FoodListingUpdate struct {
	Title          *string
	Description    *string
	Latitude       *float64
	Longitude      *float64
	FoodType       *FoodCategory
	Quantity       *int32
	AvailableFrom  *time.Time
	AvailableUntil *time.Time
}

type FoodListingRepository interface {
	InsertNewFoodListing(ctx context.Context, foods *FoodListing) (*FoodListing, error)
	GetActiveFoodListing(ctx context.Context, status string) (*[]FoodListing, error)
	GetFoodListingByID(ctx context.Context, id uuid.UUID) (*FoodListing, error)
	UpdateFoodListing(ctx context.Context, id uuid.UUID, update FoodListingUpdate) (*FoodListing, error)
	CancelFoodListing(ctx context.Context, id uuid.UUID) (*FoodListing, error)
	ReserveFoodQuantity(ctx context.Context, listingID uuid.UUID, quantity int32) (*FoodListing, error)
	ReleaseFoodQuantity(ctx context.Context, listingID uuid.UUID, quantity int32) (*FoodListing, error)
	UpdateExpiredAndClaimedListings(ctx context.Context) error
}
