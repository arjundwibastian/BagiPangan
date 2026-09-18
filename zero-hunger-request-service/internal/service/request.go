package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zero-hunger/request-service/internal/domain"
)

type RequestService struct {
	requests domain.Repository
	users    domain.UserValidator
	ttl      time.Duration
}

func NewRequestService(requests domain.Repository, users domain.UserValidator, ttl time.Duration) *RequestService {
	return &RequestService{requests: requests, users: users, ttl: ttl}
}

type CreateInput struct {
	UserID    uuid.UUID
	Latitude  float64
	Longitude float64
	RadiusKM  float64
}

type CreateResult struct {
	Request domain.FoodRequest `json:"request"`
}

func (s *RequestService) Create(ctx context.Context, in CreateInput) (CreateResult, error) {
	if in.UserID == uuid.Nil || in.Latitude < -90 || in.Latitude > 90 || in.Longitude < -180 || in.Longitude > 180 || in.RadiusKM <= 0 {
		return CreateResult{}, domain.ErrInvalidInput
	}

	if s.users != nil {
		if err := s.users.ValidateRecipient(ctx, in.UserID); err != nil {
			log.Printf("event=request_recipient_validation_failed user_id=%s error=%v", in.UserID, err)
			return CreateResult{}, err
		}
	}

	request := domain.FoodRequest{
		ID: uuid.New(), UserID: in.UserID, Latitude: in.Latitude, Longitude: in.Longitude,
		RadiusKM: in.RadiusKM, Status: domain.StatusSearching, ExpiresAt: time.Now().Add(s.ttl),
	}

	created, err := s.requests.Create(ctx, request)
	if err != nil {
		log.Printf("event=request_create_failed user_id=%s error=%v", in.UserID, err)
		return CreateResult{}, err
	}
	log.Printf("event=request_created request_id=%s user_id=%s status=%s expires_at=%s", created.ID, created.UserID, created.Status, created.ExpiresAt.Format(time.RFC3339))

	return CreateResult{Request: created}, nil
}

func (s *RequestService) Get(ctx context.Context, id uuid.UUID) (domain.FoodRequest, error) {
	if id == uuid.Nil {
		return domain.FoodRequest{}, domain.ErrInvalidInput
	}

	r, err := s.requests.GetByID(ctx, id)
	if err != nil {
		return r, err
	}

	return s.expireIfNeeded(ctx, r)
}

func (s *RequestService) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.FoodRequest, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	requests, err := s.requests.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range requests {
		requests[i], err = s.expireIfNeeded(ctx, requests[i])
		if err != nil {
			return nil, err
		}
	}
	return requests, nil
}

func (s *RequestService) Cancel(ctx context.Context, id, userID uuid.UUID) (domain.FoodRequest, error) {
	if id == uuid.Nil || userID == uuid.Nil {
		return domain.FoodRequest{}, domain.ErrInvalidInput
	}

	r, err := s.Get(ctx, id)
	if err != nil {
		return domain.FoodRequest{}, err
	}

	if r.UserID != userID {
		return domain.FoodRequest{}, domain.ErrForbidden
	}

	if r.Status != domain.StatusSearching {
		return domain.FoodRequest{}, domain.ErrConflict
	}

	return s.requests.Cancel(ctx, id, userID)
}

func (s *RequestService) MarkClaimed(ctx context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	if id == uuid.Nil || claimID == uuid.Nil {
		return domain.FoodRequest{}, domain.ErrInvalidInput
	}

	r, err := s.Get(ctx, id)
	if err != nil {
		return domain.FoodRequest{}, err
	}

	if r.Status != domain.StatusSearching {
		return domain.FoodRequest{}, domain.ErrConflict
	}

	request, err := s.requests.MarkClaimed(ctx, id, claimID)
	if err == nil {
		log.Printf("event=request_claimed request_id=%s claim_id=%s", id, claimID)
	}
	return request, err
}

func (s *RequestService) MarkCompleted(ctx context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	if id == uuid.Nil || claimID == uuid.Nil {
		return domain.FoodRequest{}, domain.ErrInvalidInput
	}

	r, err := s.Get(ctx, id)
	if err != nil {
		return domain.FoodRequest{}, err
	}

	if r.Status != domain.StatusClaimed || r.ClaimID == nil || *r.ClaimID != claimID {
		return domain.FoodRequest{}, domain.ErrConflict
	}

	request, err := s.requests.MarkCompleted(ctx, id, claimID)
	if err == nil {
		log.Printf("event=request_completed request_id=%s claim_id=%s", id, claimID)
	}
	return request, err
}

func (s *RequestService) MarkSearching(ctx context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	if id == uuid.Nil || claimID == uuid.Nil {
		return domain.FoodRequest{}, domain.ErrInvalidInput
	}
	r, err := s.Get(ctx, id)
	if err != nil {
		return domain.FoodRequest{}, err
	}
	if r.Status != domain.StatusClaimed || r.ClaimID == nil || *r.ClaimID != claimID {
		return domain.FoodRequest{}, domain.ErrConflict
	}
	request, err := s.requests.MarkSearching(ctx, id, claimID)
	if err == nil {
		log.Printf("event=request_reopened request_id=%s claim_id=%s", id, claimID)
	}
	return request, err
}

func (s *RequestService) expireIfNeeded(ctx context.Context, r domain.FoodRequest) (domain.FoodRequest, error) {
	if r.Status == domain.StatusSearching && !r.ExpiresAt.IsZero() && time.Now().After(r.ExpiresAt) {
		if err := s.requests.MarkExpired(ctx, r.ID); err != nil {
			return domain.FoodRequest{}, err
		}
		r.Status = domain.StatusExpired
		log.Printf("event=request_expired request_id=%s", r.ID)
	}

	return r, nil
}

func ParseStatus(value string) (domain.RequestStatus, error) {
	status := domain.RequestStatus(strings.ToLower(strings.TrimSpace(value)))
	switch status {
	case domain.StatusSearching, domain.StatusClaimed, domain.StatusCompleted, domain.StatusCancelled, domain.StatusExpired:
		return status, nil
	default:
		return "", errors.New("invalid request status")
	}
}
