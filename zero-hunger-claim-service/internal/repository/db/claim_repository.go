package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zero-hunger/claim-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type claimRepository struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) domain.ClaimRepository {
	return &claimRepository{db: db}
}

func (r *claimRepository) Create(
	ctx context.Context,
	claim *domain.FoodClaim,
) (*domain.FoodClaim, error) {
	if err := r.db.WithContext(ctx).Create(claim).Error; err != nil {
		return nil, err
	}

	return claim, nil
}

func (r *claimRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.FoodClaim, error) {
	var claim domain.FoodClaim

	err := r.db.
		WithContext(ctx).
		Where("claim_id = ?", id).
		First(&claim).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func (r *claimRepository) DeleteByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	result := r.db.
		WithContext(ctx).
		Delete(&domain.FoodClaim{}, "claim_id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *claimRepository) VerifyPickupCode(
	ctx context.Context,
	id uuid.UUID,
	code string,
) (*domain.FoodClaim, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var claim domain.FoodClaim

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("claim_id = ?", id).
		First(&claim).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, domain.ErrNotFound
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if claim.Status != domain.StatusWaitingForPickup {
		tx.Rollback()
		return nil, domain.ErrClaimAlreadyVerified
	}

	if time.Now().After(claim.CodeExpiresAt) {
		tx.Rollback()
		return nil, domain.ErrClaimExpired
	}

	if claim.ClaimCode != code {
		tx.Rollback()
		return nil, domain.ErrInvalidInput
	}

	now := time.Now()

	err = tx.Model(&claim).Updates(map[string]interface{}{
		"status":     domain.StatusPickedUp,
		"claimed_at": now,
		"updated_at": now,
	}).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	claim.Status = domain.StatusPickedUp
	claim.ClaimedAt = &now
	claim.UpdatedAt = now

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &claim, nil
}

func (r *claimRepository) Cancel(
	ctx context.Context,
	id uuid.UUID,
) (*domain.FoodClaim, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var claim domain.FoodClaim

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("claim_id = ?", id).
		First(&claim).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, domain.ErrNotFound
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if claim.Status != domain.StatusWaitingForPickup {
		tx.Rollback()
		return nil, domain.ErrClaimAlreadyVerified
	}

	now := time.Now()

	err = tx.Model(&claim).Updates(map[string]interface{}{
		"status":     domain.StatusCancelled,
		"updated_at": now,
	}).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	claim.Status = domain.StatusCancelled
	claim.UpdatedAt = now

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &claim, nil
}
