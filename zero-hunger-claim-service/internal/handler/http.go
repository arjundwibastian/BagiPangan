package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/claim-service/internal/domain"
	"github.com/zero-hunger/claim-service/internal/dto"
)

type HTTPHandler struct {
	claims domain.ClaimUseCase
}

func NewHTTPHandler(claims domain.ClaimUseCase) *HTTPHandler {
	return &HTTPHandler{claims: claims}
}

// CreateClaim godoc
// @Summary Create a food claim
// @Description Claim food from a listing for a specific food request
// @ID create-claim
// @Tags Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateClaimRequest true "Claim information"
// @Success 201 {object} dto.HandlerResponse{}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 403 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 409 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /claims [post]
func (h *HTTPHandler) CreateClaim(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return writeError(c, err)
	}

	var input dto.CreateClaimRequest

	if err := c.Bind(&input); err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	requestID, err := uuid.Parse(input.RequestID)
	if err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	foodListingID, err := uuid.Parse(input.FoodListingID)
	if err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		10*time.Second,
	)
	defer cancel()

	claim, err := h.claims.CreateClaim(
		ctx,
		userID,
		domain.CreateClaimInput{
			RequestID:       requestID,
			FoodListingID:   foodListingID,
			ClaimedQuantity: input.ClaimedQuantity,
		},
	)
	if err != nil {
		return writeError(c, err)
	}
	_ = claim
	return c.JSON(http.StatusCreated, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "claim created successfully",
		ResponseData:    claim,
	})
}

// GetClaim godoc
// @Summary Get claim by ID
// @Description Return a claim belonging to the authenticated user
// @ID get-claim
// @Tags Claims
// @Produce json
// @Security BearerAuth
// @Param claim_id path string true "Claim ID" format(uuid)
// @Success 200 {object} dto.HandlerResponse{responseData=domain.FoodClaim}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 403 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /claims/{claim_id} [get]
func (h *HTTPHandler) GetClaim(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return writeError(c, err)
	}

	claimID, err := uuid.Parse(c.Param("claim_id"))
	if err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	claim, err := h.claims.GetClaim(
		c.Request().Context(),
		userID,
		claimID,
	)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "claim retrieved successfully",
		ResponseData:    claim,
	})
}

// VerifyPickupCode godoc
// @Summary Verify a pickup code
// @Description Verify the claim pickup code and mark the claim as completed
// @ID verify-pickup-code
// @Tags Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param claim_id path string true "Claim ID" format(uuid)
// @Param request body dto.VerifyPickupRequest true "Pickup verification information"
// @Success 200 {object} dto.HandlerResponse{responseData=domain.FoodClaim}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 403 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 409 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /claims/{claim_id}/verify-pickup [post]
func (h *HTTPHandler) VerifyPickupCode(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return writeError(c, err)
	}

	claimID, err := uuid.Parse(c.Param("claim_id"))
	if err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	var input dto.VerifyPickupRequest
	if err := c.Bind(&input); err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	claim, err := h.claims.VerifyPickupCode(
		c.Request().Context(),
		userID,
		claimID,
		input.ClaimCode,
	)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "pickup verified successfully",
		ResponseData:    claim,
	})
}

// CancelClaim godoc
// @Summary Cancel a claim
// @Description Cancel a food claim belonging to the authenticated user
// @ID cancel-claim
// @Tags Claims
// @Produce json
// @Security BearerAuth
// @Param claim_id path string true "Claim ID" format(uuid)
// @Success 200 {object} dto.HandlerResponse{responseData=domain.FoodClaim}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 403 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 409 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /claims/{claim_id}/cancel [post]
func (h *HTTPHandler) CancelClaim(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return writeError(c, err)
	}

	claimID, err := uuid.Parse(c.Param("claim_id"))
	if err != nil {
		return writeError(c, domain.ErrInvalidInput)
	}

	claim, err := h.claims.CancelClaim(
		c.Request().Context(),
		userID,
		claimID,
	)
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "claim cancelled successfully",
		ResponseData:    claim,
	})
}

func currentUserID(c echo.Context) (uuid.UUID, error) {
	value := c.Get("user_id")

	userID, ok := value.(string)
	if !ok {
		return uuid.Nil, domain.ErrForbidden
	}

	return uuid.Parse(userID)
}

func writeError(c echo.Context, err error) error {
	statusCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest

	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound

	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden

	case errors.Is(err, domain.ErrClaimExpired),
		errors.Is(err, domain.ErrClaimAlreadyVerified),
		errors.Is(err, domain.ErrInsufficientQuantity):
		statusCode = http.StatusConflict
	}

	return c.JSON(statusCode, dto.HandlerResponse{
		ResponseCode:    "01",
		ResponseMessage: err.Error(),
		ResponseData:    nil,
	})
}
