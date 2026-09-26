package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zero-hunger/request-service/internal/domain"
)

type requestRepoMock struct {
	requests map[uuid.UUID]domain.FoodRequest
	expired  int
}

func (m *requestRepoMock) Create(_ context.Context, r domain.FoodRequest) (domain.FoodRequest, error) {
	r.CreatedAt = time.Now()
	r.UpdatedAt = r.CreatedAt
	m.requests[r.ID] = r
	return r, nil
}
func (m *requestRepoMock) GetByID(_ context.Context, id uuid.UUID) (domain.FoodRequest, error) {
	r, ok := m.requests[id]
	if !ok {
		return domain.FoodRequest{}, domain.ErrNotFound
	}
	return r, nil
}
func (m *requestRepoMock) ListByUser(_ context.Context, userID uuid.UUID) ([]domain.FoodRequest, error) {
	var out []domain.FoodRequest
	for _, r := range m.requests {
		if r.UserID == userID {
			out = append(out, r)
		}
	}
	return out, nil
}
func (m *requestRepoMock) Cancel(_ context.Context, id, userID uuid.UUID) (domain.FoodRequest, error) {
	r, err := m.GetByID(context.Background(), id)
	if err != nil {
		return r, err
	}
	if r.UserID != userID {
		return domain.FoodRequest{}, domain.ErrNotFound
	}
	r.Status = domain.StatusCancelled
	m.requests[id] = r
	return r, nil
}
func (m *requestRepoMock) MarkClaimed(_ context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	r, err := m.GetByID(context.Background(), id)
	if err != nil {
		return r, err
	}
	r.Status = domain.StatusClaimed
	r.ClaimID = &claimID
	m.requests[id] = r
	return r, nil
}
func (m *requestRepoMock) MarkSearching(_ context.Context, id, claimID uuid.UUID) (domain.FoodRequest, error) {
	r, err := m.GetByID(context.Background(), id)
	if err != nil || r.ClaimID == nil || *r.ClaimID != claimID {
		return domain.FoodRequest{}, domain.ErrConflict
	}
	r.Status = domain.StatusSearching
	r.ClaimID = nil
	m.requests[id] = r
	return r, nil
}
func (m *requestRepoMock) MarkCompleted(_ context.Context, id, _ uuid.UUID) (domain.FoodRequest, error) {
	r, err := m.GetByID(context.Background(), id)
	if err != nil {
		return r, err
	}
	r.Status = domain.StatusCompleted
	m.requests[id] = r
	return r, nil
}
func (m *requestRepoMock) MarkExpired(_ context.Context, id uuid.UUID) error {
	r, err := m.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	r.Status = domain.StatusExpired
	m.requests[id] = r
	m.expired++
	return nil
}

func TestCreateRequest(t *testing.T) {
	repo := &requestRepoMock{requests: map[uuid.UUID]domain.FoodRequest{}}
	svc := NewRequestService(repo, nil, 24*time.Hour)
	userID := uuid.New()
	out, err := svc.Create(context.Background(), CreateInput{UserID: userID, Latitude: -6.2, Longitude: 106.8, RadiusKM: 5})
	if err != nil {
		t.Fatal(err)
	}
	if out.Request.Status != domain.StatusSearching {
		t.Fatalf("status = %s", out.Request.Status)
	}
	if out.Request.ID == uuid.Nil {
		t.Fatal("request ID should be generated")
	}
}

func TestClaimAndCompleteLifecycle(t *testing.T) {
	repo := &requestRepoMock{requests: map[uuid.UUID]domain.FoodRequest{}}
	svc := NewRequestService(repo, nil, time.Hour)
	r, err := svc.Create(context.Background(), CreateInput{UserID: uuid.New(), Latitude: 0, Longitude: 0, RadiusKM: 1})
	if err != nil {
		t.Fatal(err)
	}
	claimID := uuid.New()
	if _, err = svc.MarkClaimed(context.Background(), r.Request.ID, claimID); err != nil {
		t.Fatal(err)
	}
	reopened, err := svc.MarkSearching(context.Background(), r.Request.ID, claimID)
	if err != nil {
		t.Fatal(err)
	}
	if reopened.Status != domain.StatusSearching || reopened.ClaimID != nil {
		t.Fatalf("reopened request = %+v", reopened)
	}
	if _, err = svc.MarkClaimed(context.Background(), r.Request.ID, claimID); err != nil {
		t.Fatal(err)
	}
	completed, err := svc.MarkCompleted(context.Background(), r.Request.ID, claimID)
	if err != nil {
		t.Fatal(err)
	}
	if completed.Status != domain.StatusCompleted {
		t.Fatalf("status = %s", completed.Status)
	}
	if _, err = svc.Cancel(context.Background(), r.Request.ID, r.Request.UserID); err != domain.ErrConflict {
		t.Fatalf("cancel after completion error = %v", err)
	}
}

func TestExpiredRequestIsMarkedLazily(t *testing.T) {
	repo := &requestRepoMock{requests: map[uuid.UUID]domain.FoodRequest{}}
	svc := NewRequestService(repo, nil, -time.Second)
	r, err := svc.Create(context.Background(), CreateInput{UserID: uuid.New(), Latitude: 0, Longitude: 0, RadiusKM: 1})
	if err != nil {
		t.Fatal(err)
	}
	expired, err := svc.Get(context.Background(), r.Request.ID)
	if err != nil {
		t.Fatal(err)
	}
	if expired.Status != domain.StatusExpired || repo.expired != 1 {
		t.Fatalf("status=%s expired_calls=%d", expired.Status, repo.expired)
	}
}

func TestCreateRejectsInvalidCoordinates(t *testing.T) {
	svc := NewRequestService(&requestRepoMock{requests: map[uuid.UUID]domain.FoodRequest{}}, nil, time.Hour)
	if _, err := svc.Create(context.Background(), CreateInput{UserID: uuid.New(), Latitude: 91, Longitude: 0, RadiusKM: 1}); err != domain.ErrInvalidInput {
		t.Fatalf("error = %v", err)
	}
}
