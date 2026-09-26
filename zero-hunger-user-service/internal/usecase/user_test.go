package usecase

import (
	"context"
	"github.com/google/uuid"
	"testing"
	"time"
	"github.com/zero-hunger/user-service/internal/domain"
)

type userRepoMock struct {
	byEmail map[string]domain.User
	byID    map[uuid.UUID]domain.User
}

func (m *userRepoMock) Create(_ context.Context, u domain.User) (domain.User, error) {
	u.ID = uuid.New()
	m.byEmail[u.Email] = u
	m.byID[u.ID] = u
	return u, nil
}
func (m *userRepoMock) GetByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}
func (m *userRepoMock) GetByEmail(_ context.Context, email string) (domain.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}
func (m *userRepoMock) UpdateProfile(_ context.Context, id uuid.UUID, name, phone string) (domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	u.Name = name
	u.Phone = phone
	m.byID[id] = u
	return u, nil
}

type tokenRepoMock struct{}

func (tokenRepoMock) Create(context.Context, domain.RefreshToken) error { return nil }
func (tokenRepoMock) GetActiveByHash(context.Context, string) (domain.RefreshToken, error) {
	return domain.RefreshToken{}, domain.ErrNotFound
}
func (tokenRepoMock) Revoke(context.Context, uuid.UUID) error { return nil }
func TestRegister(t *testing.T) {
	r := &userRepoMock{byEmail: map[string]domain.User{}, byID: map[uuid.UUID]domain.User{}}
	s := NewUserService(r, tokenRepoMock{}, "secret", time.Hour, time.Hour)
	u, err := s.Register(context.Background(), RegisterInput{Name: "John", Email: "john@example.com", Phone: "08123", Password: "password", Role: domain.RoleDonor})
	if err != nil {
		t.Fatal(err)
	}
	if u.PasswordHash != "" {
		t.Fatal("password hash should not be exposed")
	}
	if len(r.byEmail["john@example.com"].PasswordHash) == 0 {
		t.Fatal("password should be hashed")
	}
}
func TestRegisterRejectsInvalidInput(t *testing.T) {
	r := &userRepoMock{byEmail: map[string]domain.User{}, byID: map[uuid.UUID]domain.User{}}
	s := NewUserService(r, tokenRepoMock{}, "secret", time.Hour, time.Hour)
	if _, err := s.Register(context.Background(), RegisterInput{Name: "", Email: "bad", Phone: "", Password: "short", Role: domain.RoleDonor}); err != domain.ErrInvalidInput {
		t.Fatalf("got %v", err)
	}
}

func TestRegisterRejectsAdminRole(t *testing.T) {
	r := &userRepoMock{byEmail: map[string]domain.User{}, byID: map[uuid.UUID]domain.User{}}
	s := NewUserService(r, tokenRepoMock{}, "secret", time.Hour, time.Hour)
	if _, err := s.Register(context.Background(), RegisterInput{Name: "Admin", Email: "admin@example.com", Phone: "08123", Password: "password", Role: domain.RoleAdmin}); err != domain.ErrInvalidInput {
		t.Fatalf("got %v", err)
	}
}
