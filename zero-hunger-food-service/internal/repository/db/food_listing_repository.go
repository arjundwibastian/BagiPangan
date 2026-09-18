package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zero-hunger/food-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormFoodListingRepo struct {
	db *gorm.DB
}

func NewFoodListingRepo(db *gorm.DB) domain.FoodListingRepository {
	return &gormFoodListingRepo{db: db}
}

func (r *gormFoodListingRepo) InsertNewFoodListing(ctx context.Context, foods *domain.FoodListing) (*domain.FoodListing, error) {
	err := r.db.WithContext(ctx).Create(foods).Error
	if err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *gormFoodListingRepo) GetActiveFoodListing(ctx context.Context, status string) (*[]domain.FoodListing, error) {
	var foodListing []domain.FoodListing
	now := time.Now()
	err := r.db.Model(&foodListing).Where("status = ? AND available_from <= ? AND available_until > ? AND quantity > 0", status, now, now).Find(&foodListing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}
	return &foodListing, nil
}

func (r *gormFoodListingRepo) GetFoodListingByID(ctx context.Context, id uuid.UUID) (*domain.FoodListing, error) {
	var listing domain.FoodListing

	err := r.db.
		WithContext(ctx).
		Where("food_listing_id = ?", id).
		First(&listing).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *gormFoodListingRepo) ReserveFoodQuantity(
	ctx context.Context,
	listingID uuid.UUID,
	quantity int32,
) (*domain.FoodListing, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()

	var listing domain.FoodListing

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("food_listing_id = ?", listingID).
		First(&listing).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, domain.ErrNotFound
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	now := time.Now()

	if listing.Status != domain.StatusAvailable ||
		now.Before(listing.AvailableFrom) ||
		!now.Before(listing.AvailableUntil) {
		tx.Rollback()
		return nil, domain.ErrListingUnavailable
	}

	if quantity <= 0 || listing.Quantity < quantity {
		tx.Rollback()
		return nil, domain.ErrInsufficientQuantity
	}

	remainingQuantity := listing.Quantity - quantity
	newStatus := domain.StatusAvailable

	if remainingQuantity == 0 {
		newStatus = domain.StatusClaimed
	}

	err = tx.Model(&listing).Updates(map[string]interface{}{
		"quantity":   remainingQuantity,
		"status":     newStatus,
		"updated_at": now,
	}).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	listing.Quantity = remainingQuantity
	listing.Status = newStatus
	listing.UpdatedAt = now

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *gormFoodListingRepo) ReleaseFoodQuantity(
	ctx context.Context,
	listingID uuid.UUID,
	quantity int32,
) (*domain.FoodListing, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()

	var listing domain.FoodListing

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("food_listing_id = ?", listingID).
		First(&listing).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, domain.ErrNotFound
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if quantity <= 0 {
		tx.Rollback()
		return nil, domain.ErrInvalidInput
	}

	newQuantity := listing.Quantity + quantity
	newStatus := domain.StatusAvailable

	if time.Now().After(listing.AvailableUntil) {
		newStatus = domain.StatusExpired
	}

	err = tx.Model(&listing).Updates(map[string]interface{}{
		"quantity":   newQuantity,
		"status":     newStatus,
		"updated_at": time.Now(),
	}).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	listing.Quantity = newQuantity
	listing.Status = newStatus

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *gormFoodListingRepo) UpdateExpiredAndClaimedListings(ctx context.Context) error {
	result := r.db.
		WithContext(ctx).
		Model(&domain.FoodListing{}).
		Where("status = ?", domain.StatusAvailable).
		Where(
			"quantity = ? OR available_until <= CURRENT_TIMESTAMP",
			0,
		).
		Updates(map[string]interface{}{
			"status": gorm.Expr(`
				CASE
					WHEN quantity = ?
						THEN ?::listing_status
					WHEN available_until <= CURRENT_TIMESTAMP
						THEN ?::listing_status
					ELSE status
				END
			`,
				0,
				string(domain.StatusClaimed),
				string(domain.StatusExpired),
			),
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		})

	return result.Error

}
