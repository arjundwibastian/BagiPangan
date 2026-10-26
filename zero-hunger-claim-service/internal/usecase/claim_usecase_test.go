package usecase

import (
	"context"
	"errors"
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

func (r *claimRepoMock) GetByListingAndCode(_ context.Context, listingID uuid.UUID, code string) (*domain.FoodClaim, error) {
	for _, claim := range r.claims {
		if claim.FoodListingID == listingID && claim.ClaimCode == code {
			return claim, nil
		}
	}
	return nil, domain.ErrInvalidInput
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
	if claim.Status != domain.StatusWaitingForPickup {
		return nil, domain.ErrClaimAlreadyVerified
	}
	if time.Now().After(claim.CodeExpiresAt) {
		return nil, domain.ErrClaimExpired
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

func (r *claimRepoMock) ListByRequestIDs(
	_ context.Context,
	requestIDs []uuid.UUID,
) ([]domain.FoodClaim, error) {
	claims := make([]domain.FoodClaim, 0)
	for _, requestID := range requestIDs {
		for _, claim := range r.claims {
			if claim.RequestID == requestID {
				claims = append(claims, *claim)
			}
		}
	}

	return claims, nil
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
	data         *requestv1.FoodRequest
	markedClaim  bool
	markedSearch bool
	completed    bool
}

func (m *requestMock) GetRequest(
	_ context.Context,
	_ string,
) (*requestv1.GetRequestResponse, error) {
	return &requestv1.GetRequestResponse{
		Request: m.data,
	}, nil
}

func (m *requestMock) ListRequestsByUser(
	_ context.Context,
	_ string,
) (*requestv1.ListRequestsByUserResponse, error) {
	return &requestv1.ListRequestsByUserResponse{
		Requests: []*requestv1.FoodRequest{m.data},
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

func TestVerifyListingPickupCode(t *testing.T) {
	for _, tc := range []struct {
		name         string
		code         string
		role         userv1.UserRole
		otherOwner   bool
		otherListing bool
		expired      bool
		status       domain.ClaimStatus
		wantErr      error
	}{
		{name: "owner verifies leading-zero code", code: "012345", role: userv1.UserRole_USER_ROLE_DONOR},
		{name: "wrong code", code: "999999", role: userv1.UserRole_USER_ROLE_DONOR, wantErr: domain.ErrInvalidInput},
		{name: "malformed code", code: "01234x", role: userv1.UserRole_USER_ROLE_DONOR, wantErr: domain.ErrInvalidInput},
		{name: "other donor", code: "012345", role: userv1.UserRole_USER_ROLE_DONOR, otherOwner: true, wantErr: domain.ErrForbidden},
		{name: "recipient cannot verify", code: "012345", role: userv1.UserRole_USER_ROLE_RECIPIENT, wantErr: domain.ErrForbidden},
		{name: "code for another listing", code: "012345", role: userv1.UserRole_USER_ROLE_DONOR, otherListing: true, wantErr: domain.ErrInvalidInput},
		{name: "expired code", code: "012345", role: userv1.UserRole_USER_ROLE_DONOR, expired: true, wantErr: domain.ErrClaimExpired},
		{name: "already collected", code: "012345", role: userv1.UserRole_USER_ROLE_DONOR, status: domain.StatusPickedUp, wantErr: domain.ErrClaimAlreadyVerified},
		{name: "cancelled", code: "012345", role: userv1.UserRole_USER_ROLE_DONOR, status: domain.StatusCancelled, wantErr: domain.ErrClaimAlreadyVerified},
	} {
		t.Run(tc.name, func(t *testing.T) {
			userID, listingID, claimID := uuid.New(), uuid.New(), uuid.New()
			claim := &domain.FoodClaim{ClaimID: claimID, RequestID: uuid.New(), FoodListingID: listingID,
				ClaimCode: "012345", ClaimedQuantity: 2, Status: domain.StatusWaitingForPickup, CodeExpiresAt: time.Now().Add(time.Hour)}
			if tc.expired {
				claim.CodeExpiresAt = time.Now().Add(-time.Minute)
			}
			if tc.status != "" {
				claim.Status = tc.status
			}
			if tc.otherListing {
				claim.FoodListingID = uuid.New()
			}
			ownerID := userID
			if tc.otherOwner {
				ownerID = uuid.New()
			}
			repo := &claimRepoMock{claims: map[uuid.UUID]*domain.FoodClaim{claimID: claim}}
			requests := &requestMock{data: &requestv1.FoodRequest{RequestId: claim.RequestID.String()}}
			service := NewClaimUseCase(repo, &claimUserMock{user: &userv1.GetUserResponse{Role: tc.role}}, requests,
				&claimFoodMock{listing: &foodv1.FoodListing{FoodListingId: listingID.String(), UserId: ownerID.String()}}, &mailMock{})
			result, err := service.VerifyListingPickupCode(context.Background(), userID, listingID, tc.code)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil {
				if result == nil || result.ClaimID != claimID || result.Status != domain.StatusPickedUp || !requests.completed {
					t.Fatal("matching claim and request must be completed")
				}
			} else if requests.completed {
				t.Fatal("rejected pickup must not complete the request")
			}
		})
	}
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
