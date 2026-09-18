package domain

import (
	"context"

	"github.com/google/uuid"
)

type FoodListingRepository interface {
	InsertNewFoodListing(ctx context.Context, foods *FoodListing) (*FoodListing, error)
	GetActiveFoodListing(ctx context.Context, status string) (*[]FoodListing, error)
	GetFoodListingByID(ctx context.Context, id uuid.UUID) (*FoodListing, error)
	ReserveFoodQuantity(ctx context.Context, listingID uuid.UUID, quantity int32) (*FoodListing, error)
	ReleaseFoodQuantity(ctx context.Context, listingID uuid.UUID, quantity int32) (*FoodListing, error)
	UpdateExpiredAndClaimedListings(ctx context.Context) error
}
