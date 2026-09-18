package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type UserRepository interface {
	Create(context.Context, User) (User, error)
	GetByID(context.Context, uuid.UUID) (User, error)
	GetByEmail(context.Context, string) (User, error)
	UpdateProfile(context.Context, uuid.UUID, string, string) (User, error)
}

type RefreshTokenRepository interface {
	Create(context.Context, RefreshToken) error
	GetActiveByHash(context.Context, string) (RefreshToken, error)
	Revoke(context.Context, uuid.UUID) error
}
