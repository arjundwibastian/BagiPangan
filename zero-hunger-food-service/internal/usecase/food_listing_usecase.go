package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jftuga/geodist"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/food-service/internal/domain"
	"github.com/zero-hunger/food-service/internal/dto"
)

type foodListingUseCase struct {
	foodListingRepo domain.FoodListingRepository
	userService     domain.UserService
	requestService  domain.RequestService
}

func NewfoodListingUseCase(foodListingRepo domain.FoodListingRepository, userService domain.UserService, requestService domain.RequestService) domain.FoodListingUseCase {
	return &foodListingUseCase{foodListingRepo: foodListingRepo, userService: userService, requestService: requestService}
}

func (uc *foodListingUseCase) NewFoodListing(ctx context.Context, userID uuid.UUID, req dto.CreateFoodListingRequest) (*domain.FoodListing, error) {
	user, err := uc.userService.GetUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	if user.GetRole() != userv1.UserRole_USER_ROLE_DONOR {
		return nil, domain.ErrForbidden
	}

	if req.AvailableUntil.Before(time.Now()) {
		return nil, domain.ErrInvalidInput
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, domain.ErrFailedToGenerateUUID
	}
	// location := geom.NewPointFlat(
	// 			geom.XY,
	// 			[]float64{
	// 			req.Longitude,
	// 			req.Latitude,
	// 			},
	// 			).SetSRID(4326)
	foodType := domain.FoodCategory(req.FoodType)
	switch foodType {
	case domain.FoodCooked, domain.FoodRaw, domain.FoodBakery:
	default:
		return nil, domain.ErrInvalidInput
	}
	foodList := &domain.FoodListing{
		FoodListingID:  id,
		UserID:         userID,
		Title:          req.Title,
		Description:    req.Description,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		FoodType:       foodType,
		Quantity:       req.Quantity,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
		Status:         domain.StatusAvailable,
	}
	newFoodListing, err := uc.foodListingRepo.InsertNewFoodListing(ctx, foodList)
	if err != nil {
		return nil, err
	}
	return newFoodListing, nil
}

func (uc *foodListingUseCase) GetActiveFoodListing(ctx context.Context) (*[]domain.FoodListing, error) {
	status := string(domain.StatusAvailable)
	foodListing, err := uc.foodListingRepo.GetActiveFoodListing(ctx, status)
	if err != nil {
		return nil, err
	}
	return foodListing, nil
}

func (uc *foodListingUseCase) SearchNearbyFoodListing(ctx context.Context, requestID uuid.UUID) (*[]domain.FoodListing, error) {
	request, err := uc.requestService.GetRequest(ctx, requestID.String())
	if err != nil {
		return nil, err
	}
	return uc.SearchNearbyFoodListings(ctx, request.Request.GetLatitude(), request.Request.GetLongitude(), request.Request.GetRadiusKm())
}

func (uc *foodListingUseCase) SearchNearbyFoodListings(ctx context.Context, latitude, longitude, radiusKM float64) (*[]domain.FoodListing, error) {
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 || radiusKM <= 0 {
		return nil, domain.ErrInvalidInput
	}
	status := string(domain.StatusAvailable)
	activeListings, err := uc.foodListingRepo.GetActiveFoodListing(ctx, status)
	if err != nil {
		return nil, err
	}
	var sortedNearbyListing []domain.FoodListing
	for _, listing := range *activeListings {
		distanceBetween := distanceBetweenListings(latitude, longitude, listing.Latitude, listing.Longitude)
		if distanceBetween <= radiusKM {
			sortedNearbyListing = append(sortedNearbyListing, listing)
		}
	}
	return &sortedNearbyListing, nil
}

func distanceBetweenListings(userLat, userLon, listingLat, listingLon float64) float64 {
	userLocation := geodist.Coord{
		Lat: userLat,
		Lon: userLon,
	}

	listingLocation := geodist.Coord{
		Lat: listingLat,
		Lon: listingLon,
	}

	_, distanceKM := geodist.HaversineDistance(
		userLocation,
		listingLocation,
	)

	return distanceKM
}

func (uc *foodListingUseCase) GetFoodListingByID(ctx context.Context, id uuid.UUID) (*domain.FoodListing, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return uc.foodListingRepo.GetFoodListingByID(ctx, id)
}

func (uc *foodListingUseCase) ReserveFoodQuantity(
	ctx context.Context,
	listingID uuid.UUID,
	claimID uuid.UUID,
	quantity int32,
) (*domain.FoodListing, error) {
	if listingID == uuid.Nil ||
		claimID == uuid.Nil ||
		quantity <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return uc.foodListingRepo.ReserveFoodQuantity(
		ctx,
		listingID,
		quantity,
	)
}

func (uc *foodListingUseCase) ReleaseFoodQuantity(
	ctx context.Context,
	listingID uuid.UUID,
	claimID uuid.UUID,
	quantity int32,
) (*domain.FoodListing, error) {
	if listingID == uuid.Nil ||
		claimID == uuid.Nil ||
		quantity <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return uc.foodListingRepo.ReleaseFoodQuantity(
		ctx,
		listingID,
		quantity,
	)
}
