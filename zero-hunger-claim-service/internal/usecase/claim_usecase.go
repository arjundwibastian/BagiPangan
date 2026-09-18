package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/zero-hunger/claim-service/internal/domain"
	foodv1 "github.com/zero-hunger/contracts/gen/food/v1"
	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
)

type claimUseCase struct {
	repo           domain.ClaimRepository
	userService    domain.UserService
	requestService domain.RequestService
	foodService    domain.FoodService
	emailService domain.EmailService
}

func NewClaimUseCase(
	repo domain.ClaimRepository,
	userService domain.UserService,
	requestService domain.RequestService,
	foodService domain.FoodService,
	emailService domain.EmailService,
) domain.ClaimUseCase {
	return &claimUseCase{
		repo:           repo,
		userService:    userService,
		requestService: requestService,
		foodService:    foodService,
		emailService: emailService,
	}
}

func (uc *claimUseCase) CreateClaim(
	ctx context.Context,
	userID uuid.UUID,
	input domain.CreateClaimInput,
) (*domain.FoodClaim, error) {
	if userID == uuid.Nil ||
		input.RequestID == uuid.Nil ||
		input.FoodListingID == uuid.Nil ||
		input.ClaimedQuantity <= 0 {
		return nil, domain.ErrInvalidInput
	}

	user, err := uc.userService.GetUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	if user.GetRole() != userv1.UserRole_USER_ROLE_RECIPIENT {
		return nil, domain.ErrForbidden
	}

	requestResponse, err := uc.requestService.GetRequest(
		ctx,
		input.RequestID.String(),
	)
	if err != nil {
		return nil, err
	}

	request := requestResponse.GetRequest()
	if request == nil {
		return nil, domain.ErrNotFound
	}

	if request.GetUserId() != userID.String() {
		return nil, domain.ErrForbidden
	}

	if request.GetStatus() != requestv1.RequestStatus_REQUEST_STATUS_SEARCHING {
		return nil, domain.ErrInvalidInput
	}

	listingResponse, err := uc.foodService.GetFoodListing(
		ctx,
		input.FoodListingID.String(),
	)
	if err != nil {
		return nil, err
	}

	listing, err := findListing(
		listingResponse,
		input.FoodListingID.String(),
	)
	if err != nil {
		return nil, err
	}

	if listing.GetStatus() != foodv1.ListingStatus_LISTING_STATUS_AVAILABLE {
		return nil, domain.ErrInsufficientQuantity
	}

	if listing.GetQuantity() < input.ClaimedQuantity {
		return nil, domain.ErrInsufficientQuantity
	}

	claimID := uuid.New()

	claimCode, err := generateClaimCode()
	if err != nil {
		return nil, domain.ErrFailedToGenerateUUID
	}

	_, err = uc.foodService.ReserveFoodQuantity(
		ctx,
		input.FoodListingID.String(),
		claimID.String(),
		input.ClaimedQuantity,
	)
	if err != nil {
		return nil, err
	}

	claim := &domain.FoodClaim{
		ClaimID:         claimID,
		RequestID:       input.RequestID,
		FoodListingID:   input.FoodListingID,
		ClaimCode:       claimCode,
		CodeExpiresAt:   time.Now().Add(15 * time.Minute),
		ClaimedQuantity: input.ClaimedQuantity,
		Status:          domain.StatusWaitingForPickup,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	createdClaim, err := uc.repo.Create(ctx, claim)
	if err != nil {
		_, releaseErr := uc.foodService.ReleaseFoodQuantity(
			ctx,
			input.FoodListingID.String(),
			claimID.String(),
			input.ClaimedQuantity,
		)

		if releaseErr != nil {
			log.Printf("failed to release food after claim creation failure: %v", releaseErr)
		}

		return nil, err
	}

	_, err = uc.requestService.MarkRequestClaimed(
		ctx,
		input.RequestID.String(),
		claimID.String(),
	)
	if err != nil {
		_ = uc.repo.DeleteByID(ctx, claimID)

		_, releaseErr := uc.foodService.ReleaseFoodQuantity(
			ctx,
			input.FoodListingID.String(),
			claimID.String(),
			input.ClaimedQuantity,
		)

		if releaseErr != nil {
			log.Printf("failed to release food after request update failure: %v", releaseErr)
		}

		return nil, err
	}

	notificationData := domain.ClaimNotificationData{
		RecipientName:  user.GetName(),
		RecipientEmail: user.GetEmail(),
		FoodName:       listing.GetTitle(), 
		Quantity:       createdClaim.ClaimedQuantity,
		ClaimCode:      createdClaim.ClaimCode,
		CodeExpiresAt:  createdClaim.CodeExpiresAt,
		Status:         string(createdClaim.Status),
	}
	go func(data domain.ClaimNotificationData) {
		if err := uc.emailService.SendClaimNotification(&data); err != nil {
			log.Printf("failed to send claim notification: %v", err)
		}
	}(notificationData)	

	return createdClaim, nil
}

func (uc *claimUseCase) GetClaim(
	ctx context.Context,
	userID uuid.UUID,
	claimID uuid.UUID,
) (*domain.FoodClaim, error) {
	if userID == uuid.Nil || claimID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	claim, err := uc.repo.GetByID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	if err := uc.authorizeClaimAccess(ctx, userID, claim); err != nil {
		return nil, err
	}

	return claim, nil
}

func (uc *claimUseCase) VerifyPickupCode(
	ctx context.Context,
	userID uuid.UUID,
	claimID uuid.UUID,
	code string,
) (*domain.FoodClaim, error) {
	if userID == uuid.Nil || claimID == uuid.Nil || code == "" {
		return nil, domain.ErrInvalidInput
	}

	claim, err := uc.repo.GetByID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	listingResponse, err := uc.foodService.GetFoodListing(
		ctx,
		claim.FoodListingID.String(),
	)
	if err != nil {
		return nil, err
	}

	listing, err := findListing(
		listingResponse,
		claim.FoodListingID.String(),
	)
	if err != nil {
		return nil, err
	}

	user, err := uc.userService.GetUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	if user.GetRole() != userv1.UserRole_USER_ROLE_DONOR ||
		listing.GetUserId() != userID.String() {
		return nil, domain.ErrForbidden
	}

	pickedUpClaim, err := uc.repo.VerifyPickupCode(ctx, claimID, code)
	if err != nil {
		return nil, err
	}
	if _, err := uc.requestService.MarkRequestCompleted(ctx, claim.RequestID.String(), claimID.String()); err != nil {
		log.Printf("failed to mark request completed after pickup claim_id=%s request_id=%s: %v", claimID, claim.RequestID, err)
		return nil, err
	}
	return pickedUpClaim, nil
}

func (uc *claimUseCase) CancelClaim(
	ctx context.Context,
	userID uuid.UUID,
	claimID uuid.UUID,
) (*domain.FoodClaim, error) {
	if userID == uuid.Nil || claimID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}
	user, err := uc.userService.GetUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	if user.GetRole() != userv1.UserRole_USER_ROLE_RECIPIENT {
		return nil, domain.ErrForbidden
	}

	claim, err := uc.repo.GetByID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	requestResponse, err := uc.requestService.GetRequest(
		ctx,
		claim.RequestID.String(),
	)
	if err != nil {
		return nil, err
	}

	request := requestResponse.GetRequest()
	if request == nil {
		return nil, domain.ErrNotFound
	}

	if request.GetUserId() != userID.String() {
		return nil, domain.ErrForbidden
	}

	cancelledClaim, err := uc.repo.Cancel(ctx, claimID)
	if err != nil {
		return nil, err
	}

	_, err = uc.foodService.ReleaseFoodQuantity(
		ctx,
		claim.FoodListingID.String(),
		claim.ClaimID.String(),
		claim.ClaimedQuantity,
	)
	if err != nil {
		return nil, err
	}

	if _, err = uc.requestService.MarkRequestSearching(ctx, claim.RequestID.String(), claim.ClaimID.String()); err != nil {
		log.Printf("failed to reopen request after claim cancellation claim_id=%s request_id=%s: %v", claim.ClaimID, claim.RequestID, err)
		return nil, err
	}

	notificationData := domain.ClaimNotificationData{
		RecipientName:  user.GetName(),
		RecipientEmail: user.GetEmail(),
		FoodName : "",
		Quantity:       cancelledClaim.ClaimedQuantity,
		ClaimCode:      cancelledClaim.ClaimCode,
		CodeExpiresAt:  cancelledClaim.CodeExpiresAt,
		Status:         string(cancelledClaim.Status),
	}
	go func(data domain.ClaimNotificationData) {
		if err := uc.emailService.SendCancelNotification(&data); err != nil {
			log.Printf("failed to send cancel notification: %v", err)
		}
	}(notificationData)	

	return cancelledClaim, nil
}

func (uc *claimUseCase) authorizeClaimAccess(
	ctx context.Context,
	userID uuid.UUID,
	claim *domain.FoodClaim,
) error {
	user, err := uc.userService.GetUser(ctx, userID.String())
	if err != nil {
		return err
	}

	if user.GetRole() == userv1.UserRole_USER_ROLE_ADMIN {
		return nil
	}

	requestResponse, err := uc.requestService.GetRequest(
		ctx,
		claim.RequestID.String(),
	)
	if err != nil {
		return err
	}

	request := requestResponse.GetRequest()
	if request != nil && request.GetUserId() == userID.String() {
		return nil
	}

	listingResponse, err := uc.foodService.GetFoodListing(
		ctx,
		claim.FoodListingID.String(),
	)
	if err != nil {
		return err
	}

	listing, err := findListing(
		listingResponse,
		claim.FoodListingID.String(),
	)
	if err != nil {
		return err
	}

	if listing.GetUserId() == userID.String() {
		return nil
	}

	return domain.ErrForbidden
}

func findListing(
	response *foodv1.GetFoodListingResponse,
	listingID string,
) (*foodv1.FoodListing, error) {
	if response == nil {
		return nil, domain.ErrNotFound
	}

	for _, listing := range response.GetFoodListing() {
		if listing != nil && listing.GetFoodListingId() == listingID {
			return listing, nil
		}
	}

	return nil, domain.ErrNotFound
}

func generateClaimCode() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", number.Int64()), nil
}
