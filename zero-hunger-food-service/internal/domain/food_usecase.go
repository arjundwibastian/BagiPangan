package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/zero-hunger/food-service/internal/dto"
)

type FoodListingUseCase interface {
	NewFoodListing(ctx context.Context, userID uuid.UUID, req dto.CreateFoodListingRequest) (*FoodListing, error)
	GetActiveFoodListing(ctx context.Context) (*[]FoodListing, error)
	SearchNearbyFoodListing(ctx context.Context, requestID uuid.UUID) (*[]FoodListing, error)
	SearchNearbyFoodListings(ctx context.Context, latitude, longitude, radiusKM float64) (*[]FoodListing, error)
	GetFoodListingByID(ctx context.Context, id uuid.UUID) (*FoodListing, error)
	ReserveFoodQuantity(ctx context.Context, listingID uuid.UUID, claimID uuid.UUID, quantity int32) (*FoodListing, error)
	ReleaseFoodQuantity(ctx context.Context, listingID uuid.UUID, claimID uuid.UUID, quantity int32) (*FoodListing, error)
}
