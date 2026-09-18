package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/food-service/internal/domain"
	"github.com/zero-hunger/food-service/internal/dto"
)

type repoMock struct {
	foods   []domain.FoodListing
	inserted *domain.FoodListing
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

	input := dto.CreateFoodListingRequest{
		Title:          "Nasi Kotak",
		Description:    &description,
		Latitude:       -6.2,
		Longitude:      106.816666,
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