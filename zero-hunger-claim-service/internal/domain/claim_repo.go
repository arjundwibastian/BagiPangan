package domain

import (
	"context"

	"github.com/google/uuid"
)

type ClaimRepository interface {
	Create(ctx context.Context, claim *FoodClaim) (*FoodClaim, error)
	GetByID(ctx context.Context, id uuid.UUID) (*FoodClaim, error)
	GetByListingAndCode(ctx context.Context, listingID uuid.UUID, code string) (*FoodClaim, error)
	ListByRequestIDs(ctx context.Context, requestIDs []uuid.UUID) ([]FoodClaim, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	VerifyPickupCode(ctx context.Context, id uuid.UUID, code string) (*FoodClaim, error)
	Cancel(ctx context.Context, id uuid.UUID) (*FoodClaim, error)
}
