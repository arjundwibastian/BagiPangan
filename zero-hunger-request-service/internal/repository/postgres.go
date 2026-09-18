package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zero-hunger/request-service/internal/domain"
)

type RequestRepository struct{ db *pgxpool.Pool }

func NewRequestRepository(db *pgxpool.Pool) *RequestRepository { return &RequestRepository{db: db} }

func (r *RequestRepository) Create(ctx context.Context, request domain.FoodRequest) (domain.FoodRequest, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO food_requests(request_id,user_id,latitude,longitude,radius_km,status,expires_at)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		RETURNING request_id,user_id,latitude,longitude,radius_km,status,claim_id,expires_at,created_at,updated_at`,
		request.ID, request.UserID, request.Latitude, request.Longitude, request.RadiusKM, request.Status, request.ExpiresAt,
	).Scan(requestFields(&request)...)
	return request, err
}

func (r *RequestRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.FoodRequest, error) {
	var request domain.FoodRequest
	err := r.db.QueryRow(ctx, selectRequest+` WHERE request_id=$1`, id).Scan(requestFields(&request)...)
	return request, mapError(request, err)
}

func (r *RequestRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.FoodRequest, error) {
	rows, err := r.db.Query(ctx, selectRequest+` WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []domain.FoodRequest
	for rows.Next() {
		var request domain.FoodRequest
		if err := rows.Scan(requestFields(&request)...); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

func (r *RequestRepository) Cancel(ctx context.Context, id, userID uuid.UUID) (domain.FoodRequest, error) {
	return r.update(ctx, `
		UPDATE food_requests SET status='cancelled',updated_at=CURRENT_TIMESTAMP
		WHERE request_id=$1 AND user_id=$2 AND status='searching'
		RETURNING `+columns, id, userID)
}

func (r *RequestRepository) MarkClaimed(ctx context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	return r.update(ctx, `
		UPDATE food_requests SET status='claimed',claim_id=$2,updated_at=CURRENT_TIMESTAMP
		WHERE request_id=$1 AND status='searching' AND expires_at>CURRENT_TIMESTAMP
		RETURNING `+columns, id, claimID)
}

func (r *RequestRepository) MarkCompleted(ctx context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	return r.update(ctx, `
		UPDATE food_requests SET status='completed',updated_at=CURRENT_TIMESTAMP
		WHERE request_id=$1 AND claim_id=$2 AND status='claimed'
		RETURNING `+columns, id, claimID)
}

func (r *RequestRepository) MarkSearching(ctx context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	return r.update(ctx, `
		UPDATE food_requests SET status='searching',claim_id=NULL,updated_at=CURRENT_TIMESTAMP
		WHERE request_id=$1 AND claim_id=$2 AND status='claimed'
		RETURNING `+columns, id, claimID)
}

func (r *RequestRepository) MarkExpired(ctx context.Context, id uuid.UUID) error {
	command, err := r.db.Exec(ctx, `UPDATE food_requests SET status='expired',updated_at=CURRENT_TIMESTAMP WHERE request_id=$1 AND status='searching' AND expires_at<=CURRENT_TIMESTAMP`, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func (r *RequestRepository) update(ctx context.Context, query string, args ...any) (domain.FoodRequest, error) {
	var request domain.FoodRequest
	err := r.db.QueryRow(ctx, query, args...).Scan(requestFields(&request)...)
	return request, mapError(request, err)
}

const columns = `request_id,user_id,latitude,longitude,radius_km,status,claim_id,expires_at,created_at,updated_at`
const selectRequest = `SELECT ` + columns + ` FROM food_requests`

func requestFields(request *domain.FoodRequest) []any {
	return []any{&request.ID, &request.UserID, &request.Latitude, &request.Longitude, &request.RadiusKM, &request.Status, &request.ClaimID, &request.ExpiresAt, &request.CreatedAt, &request.UpdatedAt}
}

func mapError(request domain.FoodRequest, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}
