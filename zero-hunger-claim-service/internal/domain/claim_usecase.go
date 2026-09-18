package domain

import (
	"context"

	"github.com/google/uuid"
)

type CreateClaimInput struct {
	RequestID       uuid.UUID
	FoodListingID   uuid.UUID
	ClaimedQuantity int32
}

type ClaimUseCase interface {
	CreateClaim(
		ctx context.Context,
		userID uuid.UUID,
		input CreateClaimInput,
	) (*FoodClaim, error)

	GetClaim(
		ctx context.Context,
		userID uuid.UUID,
		claimID uuid.UUID,
	) (*FoodClaim, error)

	VerifyPickupCode(
		ctx context.Context,
		userID uuid.UUID,
		claimID uuid.UUID,
		code string,
	) (*FoodClaim, error)

	CancelClaim(
		ctx context.Context,
		userID uuid.UUID,
		claimID uuid.UUID,
	) (*FoodClaim, error)
}
