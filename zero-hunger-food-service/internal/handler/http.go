package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/food-service/internal/domain"
	"github.com/zero-hunger/food-service/internal/dto"
)

type FoodListingHandler struct {
	foodListingUC domain.FoodListingUseCase
}

func NewFoodListingHandler(FoodListingUC domain.FoodListingUseCase) *FoodListingHandler {
	return &FoodListingHandler{foodListingUC: FoodListingUC}
}

// CreateFoodListingHandler godoc
// @Summary Create a food listing
// @Description Create a new food listing for the authenticated donor
// @ID create-food-listing
// @Tags Food Listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateFoodListingRequest true "Food listing information"
// @Success 201 {object} dto.HandlerResponse{responseData=domain.FoodListing}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 403 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /food-listings [post]
func (h *FoodListingHandler) CreateFoodListingHandler(c echo.Context) error {
	userID, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}
	var req dto.CreateFoodListingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	FoodListing, err := h.foodListingUC.NewFoodListing(ctx, userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrFailedToGenerateUUID) {
			return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		} else if errors.Is(err, domain.ErrInvalidInput) {
			return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
	}
	return c.JSON(http.StatusCreated, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "successfully created FoodListing",
		ResponseData:    FoodListing,
	})
}

// GetActiveFoodListingHandler godoc
// @Summary Get active food listings
// @Description Return all currently active food listings
// @ID get-active-food-listings
// @Tags Food Listings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.HandlerResponse{responseData=[]domain.FoodListing}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /food-listings [get]
func (h *FoodListingHandler) GetActiveFoodListingHandler(c echo.Context) error {
	_, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	foods, err := h.foodListingUC.GetActiveFoodListing(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if len(*foods) == 0 {
		return c.JSON(http.StatusNotFound, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "Error: data not found",
			ResponseData:    nil,
		})
	}

	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "",
		ResponseData:    foods,
	})
}

// GetNearbyFoodListingHandler godoc
// @Summary Find nearby food listings
// @Description Find active food listings near the location of a food request
// @ID get-nearby-food-listings
// @Tags Food Listings
// @Produce json
// @Security BearerAuth
// @Param request_id query string true "Food request ID" format(uuid)
// @Success 200 {object} dto.HandlerResponse{responseData=[]domain.FoodListing}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /food-listings/nearby [get]
func (h *FoodListingHandler) GetNearbyFoodListingHandler(c echo.Context) error {
	_, err := uuid.Parse(c.Get("user_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}
	requestID := c.QueryParam("request_id")
	if requestID == "" {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "request_id is required",
			ResponseData:    nil,
		})
	}
	parsedRequestID, err := uuid.Parse(requestID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})

	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	foods, err := h.foodListingUC.SearchNearbyFoodListing(ctx, parsedRequestID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if len(*foods) == 0 {
		return c.JSON(http.StatusNotFound, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "Error: no nearby listing found",
			ResponseData:    nil,
		})
	}

	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "",
		ResponseData:    foods,
	})
}

// GetFoodListingByIDHandler godoc
// @Summary Get food listing by ID
// @Description Return the details of a specific food listing
// @ID get-food-listing-by-id
// @Tags Food Listings
// @Produce json
// @Security BearerAuth
// @Param food_listing_id path string true "Food listing ID" format(uuid)
// @Success 200 {object} dto.HandlerResponse{responseData=domain.FoodListing}
// @Failure 400 {object} dto.HandlerResponse
// @Failure 401 {object} dto.HandlerResponse
// @Failure 404 {object} dto.HandlerResponse
// @Failure 500 {object} dto.HandlerResponse
// @Router /food-listings/{food_listing_id} [get]
func (h *FoodListingHandler) GetFoodListingByIDHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("food_listing_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "invalid food listing id",
			ResponseData:    nil,
		})
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		10*time.Second,
	)
	defer cancel()

	listing, err := h.foodListingUC.GetFoodListingByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "food listing not found",
				ResponseData:    nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}

	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "successfully retrieved food listing",
		ResponseData:    listing,
	})
}
