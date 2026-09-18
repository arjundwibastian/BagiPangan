package dto

import (
	"time"
)

type CreateFoodListingRequest struct {
	Title          string    `json:"title" validate:"required,min=3,max=255"`
	Description    *string   `json:"description,omitempty" validate:"omitempty,max=2000"`
	Longitude      float64   `json:"longitude" validate:"required,gte=-180,lte=180"`
	Latitude       float64   `json:"latitude" validate:"required,gte=-90,lte=90"`
	FoodType       string    `json:"food_type" validate:"required,oneof=cooked_food raw_food bakery"`
	Quantity       int32     `json:"quantity" validate:"gte=0"`
	AvailableFrom  time.Time `json:"available_from" validate:"required"`
	AvailableUntil time.Time `json:"available_until" validate:"required,gtfield=AvailableFrom"`
}
