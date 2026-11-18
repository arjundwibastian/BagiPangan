package usecase

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/food-service/internal/domain"
	"github.com/zero-hunger/food-service/internal/dto"
)

func TestNearbyRejectsNonFiniteNumbers(t *testing.T) {
	svc := NewfoodListingUseCase(nil, nil, nil)
	for _, values := range [][3]float64{{math.NaN(), 0, 1}, {0, math.NaN(), 1}, {0, 0, math.NaN()}, {math.Inf(1), 0, 1}, {0, math.Inf(-1), 1}, {0, 0, math.Inf(1)}} {
		if _, err := svc.SearchNearbyFoodListings(context.Background(), values[0], values[1], values[2]); !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("%v: %v", values, err)
		}
	}
}

type nearbyRequestMock struct{ response *requestv1.GetRequestResponse }

func (m nearbyRequestMock) GetRequest(context.Context, string) (*requestv1.GetRequestResponse, error) {
	return m.response, nil
}

func TestNearbyRequestOwnership(t *testing.T) {
	owner := uuid.New()
	for _, tc := range []struct {
		name    string
		caller  uuid.UUID
		missing bool
		want    error
	}{
		{name: "owner", caller: owner},
		{name: "other user", caller: uuid.New(), want: domain.ErrForbidden},
		{name: "missing request", caller: owner, missing: true, want: domain.ErrNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := &requestv1.GetRequestResponse{}
			if !tc.missing {
				response.Request = &requestv1.FoodRequest{UserId: owner.String(), RadiusKm: 5}
			}
			// A nil repository makes accidental searches on denied access fail loudly.
			var repo domain.FoodListingRepository
			if tc.want == nil {
				repo = &repoMock{}
			}
			svc := NewfoodListingUseCase(repo, nil, nearbyRequestMock{response: response})
			_, err := svc.SearchNearbyFoodListing(context.Background(), tc.caller, uuid.New())
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

type repoMock struct {
	foods    []domain.FoodListing
	inserted *domain.FoodListing
}

func (r *repoMock) ListByUser(_ context.Context, userID uuid.UUID) ([]domain.FoodListing, error) {
	listings := make([]domain.FoodListing, 0)
	for _, listing := range r.foods {
		if listing.UserID == userID {
			listings = append(listings, listing)
		}
	}
	return listings, nil
}

func TestListOwnFoodListings(t *testing.T) {
	owner := uuid.New()
	repo := &repoMock{foods: []domain.FoodListing{
		{FoodListingID: uuid.New(), UserID: owner, Status: domain.StatusAvailable, Quantity: 3},
		{FoodListingID: uuid.New(), UserID: owner, Status: domain.StatusClaimed, Quantity: 0},
		{FoodListingID: uuid.New(), UserID: uuid.New(), Status: domain.StatusAvailable, Quantity: 5},
	}}
	service := NewfoodListingUseCase(repo, nil, nil)
	listings, err := service.ListOwnFoodListings(context.Background(), owner)
	if err != nil || len(listings) != 2 {
		t.Fatalf("expected both owned listings, got %v, %v", listings, err)
	}
	if listings[1].Status != domain.StatusClaimed {
		t.Fatal("fully reserved listing must remain visible")
	}
	empty, err := service.ListOwnFoodListings(context.Background(), uuid.New())
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("expected an empty list, got %v, %v", empty, err)
	}
}

func (r *repoMock) InsertNewFoodListing(
	_ context.Context,
	food *domain.FoodListing,
) (*domain.FoodListing, error) {
	r.inserted = food
	return food, nil
}

func (r *repoMock) GetActiveFoodListing(
	_ context.Context,
	_ string,
) (*[]domain.FoodListing, error) {
	return &r.foods, nil
}

func (r *repoMock) GetFoodListingByID(
	_ context.Context,
	id uuid.UUID,
) (*domain.FoodListing, error) {
	for i := range r.foods {
		if r.foods[i].FoodListingID == id {
			return &r.foods[i], nil
		}
	}

	return nil, domain.ErrNotFound
}

func (r *repoMock) ReserveFoodQuantity(
	_ context.Context,
	id uuid.UUID,
	quantity int32,
) (*domain.FoodListing, error) {
	food, err := r.GetFoodListingByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	food.Quantity -= quantity
	return food, nil
}

func (r *repoMock) ReleaseFoodQuantity(
	_ context.Context,
	id uuid.UUID,
	quantity int32,
) (*domain.FoodListing, error) {
	food, err := r.GetFoodListingByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	food.Quantity += quantity
	return food, nil
}

func (r *repoMock) UpdateExpiredAndClaimedListings(
	_ context.Context,
) error {
	return nil
}

func (r *repoMock) UpdateFoodListing(
	_ context.Context,
	id uuid.UUID,
	update domain.FoodListingUpdate,
) (*domain.FoodListing, error) {
	food, err := r.GetFoodListingByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if update.Title != nil {
		food.Title = *update.Title
	}
	if update.Quantity != nil {
		food.Quantity = *update.Quantity
	}

	return food, nil
}

func (r *repoMock) CancelFoodListing(
	_ context.Context,
	id uuid.UUID,
) (*domain.FoodListing, error) {
	food, err := r.GetFoodListingByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	food.Status = domain.StatusExpired
	return food, nil
}

type mockUser struct {
	data *userv1.GetUserResponse
}

func (m *mockUser) GetUser(
	_ context.Context,
	_ string,
) (*userv1.GetUserResponse, error) {
	return m.data, nil
}

func TestCreateFoodListing(t *testing.T) {
	userID := uuid.New()
	repository := &repoMock{}

	userClient := &mockUser{
		data: &userv1.GetUserResponse{
			UserId: userID.String(),
			Name:   "Budi",
			Email:  "budi@example.com",
			Role:   userv1.UserRole_USER_ROLE_DONOR,
		},
	}

	service := NewfoodListingUseCase(repository, userClient, nil)

	description := "Fresh food"
	now := time.Now()
	latitude, longitude := -6.2, 106.816666

	input := dto.CreateFoodListingRequest{
		Title:          "Nasi Kotak",
		Description:    &description,
		Latitude:       &latitude,
		Longitude:      &longitude,
		FoodType:       "cooked_food",
		Quantity:       10,
		AvailableFrom:  now.Add(time.Minute),
		AvailableUntil: now.Add(time.Hour),
	}

	result, err := service.NewFoodListing(
		context.Background(),
		userID,
		input,
	)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("food listing should not be nil")
	}

	if result.UserID != userID {
		t.Fatalf("expected user %s, got %s", userID, result.UserID)
	}

	if result.Title != "Nasi Kotak" {
		t.Fatalf("unexpected title: %s", result.Title)
	}

	if result.Quantity != 10 {
		t.Fatalf("unexpected quantity: %d", result.Quantity)
	}

	if result.Status != domain.StatusAvailable {
		t.Fatalf("unexpected status: %s", result.Status)
	}

	if repository.inserted == nil {
		t.Fatal("repository was not called")
	}
}

func TestNearbyFood(t *testing.T) {
	closeFood := uuid.New()

	dbMock := &repoMock{
		foods: []domain.FoodListing{
			{
				FoodListingID: closeFood,
				Title:         "Nearby Food",
				Latitude:      -6.2001,
				Longitude:     106.8167,
				Quantity:      5,
				Status:        domain.StatusAvailable,
			},
			{
				FoodListingID: uuid.New(),
				Title:         "Far Food",
				Latitude:      -7.250445,
				Longitude:     112.768845,
				Quantity:      5,
				Status:        domain.StatusAvailable,
			},
		},
	}

	foodService := NewfoodListingUseCase(dbMock, nil, nil)

	result, err := foodService.SearchNearbyFoodListings(
		context.Background(),
		-6.2,
		106.816666,
		5,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(*result) != 1 {
		t.Fatalf("expected 1 food listing, got %d", len(*result))
	}

	if (*result)[0].FoodListingID != closeFood {
		t.Fatal("returned the wrong food listing")
	}
}
