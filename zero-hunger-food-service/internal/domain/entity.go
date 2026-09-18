package domain

import (
	"time"

	"github.com/google/uuid"
)

type FoodCategory string

const (
	FoodCooked FoodCategory = "cooked_food"
	FoodRaw    FoodCategory = "raw_food"
	FoodBakery FoodCategory = "bakery"
)

type ListingStatus string

const (
	StatusAvailable ListingStatus = "available"
	StatusClaimed   ListingStatus = "claimed"
	StatusExpired   ListingStatus = "expired"
)

type FoodListing struct {
	FoodListingID  uuid.UUID     `gorm:"column:food_listing_id;type:uuid;default:gen_random_uuid();primaryKey" json:"food_listing_id"`
	UserID         uuid.UUID     `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	Title          string        `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description    *string       `gorm:"column:description;type:text" json:"description,omitempty"`
	Latitude       float64       `gorm:"column:latitude;type:decimal(9,6);not null" json:"latitude"`
	Longitude      float64       `gorm:"column:longitude;type:decimal(9,6);not null" json:"longitude"`
	FoodType       FoodCategory  `gorm:"column:food_type;type:food_category;not null;default:cooked_food" json:"food_type"`
	Quantity       int32         `gorm:"column:quantity;not null;check:quantity >= 0" json:"quantity"`
	AvailableFrom  time.Time     `gorm:"column:available_from;type:timestamptz;not null" json:"available_from"`
	AvailableUntil time.Time     `gorm:"column:available_until;type:timestamptz;not null" json:"available_until"`
	Status         ListingStatus `gorm:"column:status;type:listing_status;not null;default:available" json:"status"`
	CreatedAt      time.Time     `gorm:"column:created_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
