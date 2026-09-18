package domain

import (
	"context"

	foodv1 "github.com/zero-hunger/contracts/gen/food/v1"
	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (*userv1.GetUserResponse, error)
}

type FoodService interface {
	GetFoodListing(ctx context.Context, foodListingID string) (*foodv1.GetFoodListingResponse, error)
	ReserveFoodQuantity(ctx context.Context, foodListingID string, claimID string, quantity int32) (*foodv1.ReserveFoodQuantityResponse, error)
	ReleaseFoodQuantity(ctx context.Context, foodListingID string, claimID string, quantity int32) (*foodv1.ReleaseFoodQuantityResponse, error)
}

type RequestService interface {
	GetRequest(
		ctx context.Context,
		requestID string,
	) (*requestv1.GetRequestResponse, error)

	MarkRequestClaimed(
		ctx context.Context,
		requestID string,
		claimID string,
	) (*requestv1.MarkRequestClaimedResponse, error)

	MarkRequestSearching(
		ctx context.Context,
		requestID string,
		claimID string,
	) (*requestv1.MarkRequestSearchingResponse, error)

	MarkRequestCompleted(
		ctx context.Context,
		requestID string,
		claimID string,
	) (*requestv1.MarkRequestCompletedResponse, error)
}
