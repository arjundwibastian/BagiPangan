package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	StatusSearching RequestStatus = "searching"
	StatusClaimed   RequestStatus = "claimed"
	StatusCompleted RequestStatus = "completed"
	StatusCancelled RequestStatus = "cancelled"
	StatusExpired   RequestStatus = "expired"
)

type FoodRequest struct {
	ID        uuid.UUID     `json:"request_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	RadiusKM  float64       `json:"radius_km"`
	Status    RequestStatus `json:"status"`
	ClaimID   *uuid.UUID    `json:"claim_id,omitempty"`
	ExpiresAt time.Time     `json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Repository interface {
	Create(context.Context, FoodRequest) (FoodRequest, error)
	GetByID(context.Context, uuid.UUID) (FoodRequest, error)
	ListByUser(context.Context, uuid.UUID) ([]FoodRequest, error)
	Cancel(context.Context, uuid.UUID, uuid.UUID) (FoodRequest, error)
	MarkClaimed(context.Context, uuid.UUID, uuid.UUID) (FoodRequest, error)
	MarkSearching(context.Context, uuid.UUID, uuid.UUID) (FoodRequest, error)
	MarkCompleted(context.Context, uuid.UUID, uuid.UUID) (FoodRequest, error)
	MarkExpired(context.Context, uuid.UUID) error
}

type UserValidator interface {
	ValidateRecipient(context.Context, uuid.UUID) error
}
