package handler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	foodv1 "github.com/zero-hunger/contracts/gen/food/v1"
	"github.com/zero-hunger/food-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FoodRPCServer struct {
	foodv1.UnimplementedFoodServiceServer
	food domain.FoodListingUseCase
}

func NewFoodRpcServer(foodUC domain.FoodListingUseCase) *FoodRPCServer {
	return &FoodRPCServer{
		food: foodUC,
	}
}

func listingStatusToProto(status domain.ListingStatus) foodv1.ListingStatus {
	switch status {
	case domain.StatusAvailable:
		return foodv1.ListingStatus_LISTING_STATUS_AVAILABLE
	case domain.StatusClaimed:
		return foodv1.ListingStatus_LISTING_STATUS_CLAIMED
	case domain.StatusExpired:
		return foodv1.ListingStatus_LISTING_STATUS_EXPIRED
	default:
		return foodv1.ListingStatus_LISTING_STATUS_UNSPECIFIED
	}
}
func foodCategoryToProto(category domain.FoodCategory) foodv1.FoodCategory {
	switch category {
	case domain.FoodCooked:
		return foodv1.FoodCategory_FOOD_CATEGORY_COOKED_FOOD
	case domain.FoodRaw:
		return foodv1.FoodCategory_FOOD_CATEGORY_RAW_FOOD
	case domain.FoodBakery:
		return foodv1.FoodCategory_FOOD_CATEGORY_BAKERY
	default:
		return foodv1.FoodCategory_FOOD_CATEGORY_UNSPECIFIED
	}
}

func mapFoodToProto(food *domain.FoodListing) *foodv1.FoodListing {
	if food == nil {
		return nil
	}
	return &foodv1.FoodListing{
		FoodListingId:  food.FoodListingID.String(),
		UserId:         food.UserID.String(),
		Title:          food.Title,
		Description:    food.Description,
		Latitude:       food.Latitude,
		Longitude:      food.Longitude,
		FoodType:       foodCategoryToProto(food.FoodType),
		Quantity:       food.Quantity,
		AvailableFrom:  timestamppb.New(food.AvailableFrom),
		AvailableUntil: timestamppb.New(food.AvailableUntil),
		Status:         listingStatusToProto(food.Status),
		CreatedAt:      timestamppb.New(food.CreatedAt),
		UpdatedAt:      timestamppb.New(food.UpdatedAt),
	}
}

func (s *FoodRPCServer) GetFoodListing(ctx context.Context, in *foodv1.GetFoodListingRequest) (*foodv1.GetFoodListingResponse, error) {
	id, err := uuid.Parse(in.GetFoodListingId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid food listing id",
		)
	}

	listing, err := s.food.GetFoodListingByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"food listing not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"failed to retrieve food listing",
		)
	}

	return &foodv1.GetFoodListingResponse{
		FoodListing: []*foodv1.FoodListing{
			mapFoodToProto(listing),
		},
	}, nil
}

func (s *FoodRPCServer) SearchNearbyFoodListings(ctx context.Context, in *foodv1.SearchNearbyFoodListingsRequest) (*foodv1.SearchNearbyFoodListingsResponse, error) {
	if in.GetLocation() == nil {
		return nil, status.Error(codes.InvalidArgument, "location is required")
	}
	listings, err := s.food.SearchNearbyFoodListings(ctx, in.GetLocation().GetLatitude(), in.GetLocation().GetLongitude(), in.GetRadiusKm())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, "invalid nearby search request")
		}
		return nil, status.Error(codes.Internal, "failed to search nearby food listings")
	}
	response := &foodv1.SearchNearbyFoodListingsResponse{FoodListings: make([]*foodv1.FoodListing, 0, len(*listings))}
	for i := range *listings {
		response.FoodListings = append(response.FoodListings, mapFoodToProto(&(*listings)[i]))
	}
	return response, nil
}

func (s *FoodRPCServer) ReserveFoodQuantity(
	ctx context.Context,
	in *foodv1.ReserveFoodQuantityRequest,
) (*foodv1.ReserveFoodQuantityResponse, error) {
	listingID, err := uuid.Parse(in.GetFoodListingId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid food listing id",
		)
	}

	claimID, err := uuid.Parse(in.GetClaimId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid claim id",
		)
	}

	if in.GetQuantity() <= 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"quantity must be greater than zero",
		)
	}

	listing, err := s.food.ReserveFoodQuantity(
		ctx,
		listingID,
		claimID,
		in.GetQuantity(),
	)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"food listing not found",
			)
		}

		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(
				codes.InvalidArgument,
				"invalid reserve request",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"failed to reserve food quantity",
		)
	}

	return &foodv1.ReserveFoodQuantityResponse{
		FoodListingId:     listing.FoodListingID.String(),
		RemainingQuantity: listing.Quantity,
		Status:            listingStatusToProto(listing.Status),
	}, nil
}

func (s *FoodRPCServer) ReleaseFoodQuantity(
	ctx context.Context,
	in *foodv1.ReleaseFoodQuantityRequest,
) (*foodv1.ReleaseFoodQuantityResponse, error) {
	listingID, err := uuid.Parse(in.GetFoodListingId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid food listing id",
		)
	}

	claimID, err := uuid.Parse(in.GetClaimId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid claim id",
		)
	}

	if in.GetQuantity() <= 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"quantity must be greater than zero",
		)
	}

	listing, err := s.food.ReleaseFoodQuantity(
		ctx,
		listingID,
		claimID,
		in.GetQuantity(),
	)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"food listing not found",
			)
		}

		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(
				codes.InvalidArgument,
				"invalid release request",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"failed to release food quantity",
		)
	}

	return &foodv1.ReleaseFoodQuantityResponse{
		FoodListingId: listing.FoodListingID.String(),
		Quantity:      listing.Quantity,
		Status:        listingStatusToProto(listing.Status),
	}, nil
}
