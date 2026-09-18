package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zero-hunger/claim-service/internal/domain"
	foodv1 "github.com/zero-hunger/contracts/gen/food/v1"
	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
)

type claimRepoMock struct {
	claims map[uuid.UUID]*domain.FoodClaim
}

func (r *claimRepoMock) Create(
	_ context.Context,
	claim *domain.FoodClaim,
) (*domain.FoodClaim, error) {
	r.claims[claim.ClaimID] = claim
	return claim, nil
}

func (r *claimRepoMock) GetByID(
	_ context.Context,
	id uuid.UUID,
) (*domain.FoodClaim, error) {
	claim, ok := r.claims[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return claim, nil
}

func (r *claimRepoMock) DeleteByID(
	_ context.Context,
	id uuid.UUID,
) error {
	delete(r.claims, id)
	return nil
}

func (r *claimRepoMock) VerifyPickupCode(
	_ context.Context,
	id uuid.UUID,
	code string,
) (*domain.FoodClaim, error) {
	claim, ok := r.claims[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	if claim.ClaimCode != code {
		return nil, domain.ErrInvalidInput
	}

	now := time.Now()
	claim.Status = domain.StatusPickedUp
	claim.ClaimedAt = &now

	return claim, nil
}

func (r *claimRepoMock) Cancel(
	_ context.Context,
	id uuid.UUID,
) (*domain.FoodClaim, error) {
	claim, ok := r.claims[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	claim.Status = domain.StatusCancelled
	return claim, nil
}

type claimUserMock struct {
	user *userv1.GetUserResponse
}

func (m *claimUserMock) GetUser(
	_ context.Context,
	_ string,
) (*userv1.GetUserResponse, error) {
	return m.user, nil
}

type requestMock struct {
	data        *requestv1.FoodRequest
	markedClaim bool
	markedSearch bool
	completed   bool
}

func (m *requestMock) GetRequest(
	_ context.Context,
	_ string,
) (*requestv1.GetRequestResponse, error) {
	return &requestv1.GetRequestResponse{
		Request: m.data,
	}, nil
}

func (m *requestMock) MarkRequestClaimed(
	_ context.Context,
	_ string,
	_ string,
) (*requestv1.MarkRequestClaimedResponse, error) {
	m.markedClaim = true
	m.data.Status = requestv1.RequestStatus_REQUEST_STATUS_CLAIMED

	return &requestv1.MarkRequestClaimedResponse{
		Request: m.data,
	}, nil
}

func (m *requestMock) MarkRequestSearching(
	_ context.Context,
	_ string,
	_ string,
) (*requestv1.MarkRequestSearchingResponse, error) {
	m.markedSearch = true
	m.data.Status = requestv1.RequestStatus_REQUEST_STATUS_SEARCHING
	return &requestv1.MarkRequestSearchingResponse{
		Request: m.data,
	}, nil
}

func (m *requestMock) MarkRequestCompleted(
	_ context.Context,
	_ string,
	_ string,
) (*requestv1.MarkRequestCompletedResponse, error) {
	m.completed = true
	m.data.Status = requestv1.RequestStatus_REQUEST_STATUS_COMPLETED

	return &requestv1.MarkRequestCompletedResponse{
		Request: m.data,
	}, nil
}

type claimFoodMock struct {
	listing  *foodv1.FoodListing
	reserved int32
	released int32
}

func (m *claimFoodMock) GetFoodListing(
	_ context.Context,
	_ string,
) (*foodv1.GetFoodListingResponse, error) {
	return &foodv1.GetFoodListingResponse{
		FoodListing: []*foodv1.FoodListing{m.listing},
	}, nil
}

func (m *claimFoodMock) ReserveFoodQuantity(
	_ context.Context,
	_ string,
	_ string,
	quantity int32,
) (*foodv1.ReserveFoodQuantityResponse, error) {
	m.reserved += quantity
	m.listing.Quantity -= quantity

	return &foodv1.ReserveFoodQuantityResponse{
		FoodListingId:     m.listing.FoodListingId,
		RemainingQuantity: m.listing.Quantity,
		Status:            m.listing.Status,
	}, nil
}

func (m *claimFoodMock) ReleaseFoodQuantity(
	_ context.Context,
	_ string,
	_ string,
	quantity int32,
) (*foodv1.ReleaseFoodQuantityResponse, error) {
	m.released += quantity
	m.listing.Quantity += quantity

	return &foodv1.ReleaseFoodQuantityResponse{
		FoodListingId: m.listing.FoodListingId,
		Quantity:      m.listing.Quantity,
		Status:        m.listing.Status,
	}, nil
}

type mailMock struct{}

func (m *mailMock) SendClaimNotification(
	_ *domain.ClaimNotificationData,
) error {
	return nil
}

func (m *mailMock) SendCancelNotification(
	_ *domain.ClaimNotificationData,
) error {
	return nil
}

func TestCreateClaim(t *testing.T) {
	userID := uuid.New()
	requestID := uuid.New()
	listingID := uuid.New()

	repository := &claimRepoMock{
		claims: make(map[uuid.UUID]*domain.FoodClaim),
	}

	userClient := &claimUserMock{
		user: &userv1.GetUserResponse{
			UserId: userID.String(),
			Name:   "Budi",
			Email:  "budi@example.com",
			Role:   userv1.UserRole_USER_ROLE_RECIPIENT,
		},
	}

	reqClient := &requestMock{
		data: &requestv1.FoodRequest{
			RequestId: requestID.String(),
			UserId:    userID.String(),
			Status:    requestv1.RequestStatus_REQUEST_STATUS_SEARCHING,
		},
	}

	foodClient := &claimFoodMock{
		listing: &foodv1.FoodListing{
			FoodListingId: listingID.String(),
			UserId:        uuid.New().String(),
			Title:         "Nasi Kotak",
			Quantity:      10,
			Status:        foodv1.ListingStatus_LISTING_STATUS_AVAILABLE,
		},
	}

	service := NewClaimUseCase(
		repository,
		userClient,
		reqClient,
		foodClient,
		&mailMock{},
	)

	result, err := service.CreateClaim(
		context.Background(),
		userID,
		domain.CreateClaimInput{
			RequestID:       requestID,
			FoodListingID:   listingID,
			ClaimedQuantity: 2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("claim should not be nil")
	}

	if result.RequestID != requestID {
		t.Fatalf("unexpected request id: %s", result.RequestID)
	}

	if result.FoodListingID != listingID {
		t.Fatalf("unexpected listing id: %s", result.FoodListingID)
	}

	if result.ClaimedQuantity != 2 {
		t.Fatalf("unexpected quantity: %d", result.ClaimedQuantity)
	}

	if result.Status != domain.StatusWaitingForPickup {
		t.Fatalf("unexpected status: %s", result.Status)
	}

	if len(result.ClaimCode) != 6 {
		t.Fatalf("unexpected claim code: %s", result.ClaimCode)
	}

	if foodClient.reserved != 2 {
		t.Fatalf("expected 2 reserved, got %d", foodClient.reserved)
	}

	if !reqClient.markedClaim {
		t.Fatal("request was not marked as claimed")
	}
}

func TestCancelClaim(t *testing.T) {
	userID := uuid.New()
	requestID := uuid.New()
	listingID := uuid.New()
	claimID := uuid.New()

	existingClaim := &domain.FoodClaim{
		ClaimID:         claimID,
		RequestID:       requestID,
		FoodListingID:   listingID,
		ClaimCode:       "123456",
		ClaimedQuantity: 3,
		Status:          domain.StatusWaitingForPickup,
		CodeExpiresAt:   time.Now().Add(15 * time.Minute),
	}

	repository := &claimRepoMock{
		claims: map[uuid.UUID]*domain.FoodClaim{
			claimID: existingClaim,
		},
	}

	userClient := &claimUserMock{
		user: &userv1.GetUserResponse{
			UserId: userID.String(),
			Name:   "Budi",
			Email:  "budi@example.com",
			Role:   userv1.UserRole_USER_ROLE_RECIPIENT,
		},
	}

	reqClient := &requestMock{
		data: &requestv1.FoodRequest{
			RequestId: requestID.String(),
			UserId:    userID.String(),
			Status:    requestv1.RequestStatus_REQUEST_STATUS_CLAIMED,
		},
	}

	foodClient := &claimFoodMock{
		listing: &foodv1.FoodListing{
			FoodListingId: listingID.String(),
			Quantity:      7,
			Status:        foodv1.ListingStatus_LISTING_STATUS_AVAILABLE,
		},
	}

	service := NewClaimUseCase(
		repository,
		userClient,
		reqClient,
		foodClient,
		&mailMock{},
	)

	result, err := service.CancelClaim(
		context.Background(),
		userID,
		claimID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != domain.StatusCancelled {
		t.Fatalf("unexpected status: %s", result.Status)
	}

	if foodClient.released != 3 {
		t.Fatalf("expected 3 released, got %d", foodClient.released)
	}

	if foodClient.listing.Quantity != 10 {
		t.Fatalf(
			"expected quantity to return to 10, got %d",
			foodClient.listing.Quantity,
		)
	}

	if !reqClient.markedSearch {
	t.Fatal("request was not marked as searching")
	}


	if reqClient.data.Status != requestv1.RequestStatus_REQUEST_STATUS_SEARCHING {
		t.Fatalf(
			"expected request status searching, got %s",
			reqClient.data.Status,
		)
	}
}