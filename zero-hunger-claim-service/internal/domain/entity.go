package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClaimStatus string

const (
	StatusWaitingForPickup ClaimStatus = "waiting_for_pickup"
	StatusPickedUp         ClaimStatus = "picked_up"
	StatusCancelled        ClaimStatus = "cancelled"
)

type FoodClaim struct {
	ClaimID         uuid.UUID   `gorm:"column:claim_id;type:uuid;default:gen_random_uuid();primaryKey" json:"claim_id"`
	RequestID       uuid.UUID   `gorm:"column:request_id;type:uuid;not null" json:"request_id"`
	FoodListingID   uuid.UUID   `gorm:"column:food_listing_id;type:uuid;not null" json:"food_listing_id"`
	ClaimCode       string      `gorm:"column:claim_code;type:varchar(10);uniqueIndex;not null" json:"claim_code"`
	CodeExpiresAt   time.Time   `gorm:"column:code_expires_at;type:timestamptz;not null" json:"code_expires_at"`
	ClaimedAt       *time.Time  `gorm:"column:claimed_at;type:timestamptz" json:"claimed_at,omitempty"`
	ClaimedQuantity int32       `gorm:"column:claimed_quantity;not null;check:claimed_quantity > 0" json:"claimed_quantity"`
	Status          ClaimStatus `gorm:"column:status;type:claim_type;not null;default:waiting_for_pickup" json:"status"`
	CreatedAt       time.Time   `gorm:"column:created_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (FoodClaim) TableName() string {
	return "food_claims"
}



type Role string

const (
	RoleAdmin     Role = "admin"
	RoleDonor     Role = "donor"
	RoleRecipient Role = "recipient"
)

type User struct {
	ID           uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}




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

type ClaimNotificationData struct {
	RecipientName string
	RecipientEmail string
	FoodName string
	Quantity int32
	ClaimCode     string
	CodeExpiresAt time.Time
	Status        string
}