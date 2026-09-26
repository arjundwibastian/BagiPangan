package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/zero-hunger/user-service/internal/domain"
)

type UserService struct {
	users                 domain.UserRepository
	tokens                domain.RefreshTokenRepository
	secret                string
	accessTTL, refreshTTL time.Duration
}

func NewUserService(users domain.UserRepository, tokens domain.RefreshTokenRepository, secret string, accessTTL, refreshTTL time.Duration) *UserService {
	return &UserService{users: users, tokens: tokens, secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

type RegisterInput struct {
	Name, Email, Phone, Password string
	Role                         domain.Role
}

func (s *UserService) Register(ctx context.Context, in RegisterInput) (domain.User, error) {
	if err := validateUserInput(in.Name, in.Email, in.Phone, in.Password); err != nil {
		return domain.User{}, err
	}
	if in.Role != domain.RoleDonor && in.Role != domain.RoleRecipient {
		return domain.User{}, domain.ErrInvalidInput
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return domain.User{}, domain.ErrAlreadyExists
	} else if !errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	u, err := s.users.Create(ctx, domain.User{Name: strings.TrimSpace(in.Name), Email: email, Phone: strings.TrimSpace(in.Phone), PasswordHash: string(hash), Role: in.Role})
	u.PasswordHash = ""
	return u, err
}

type AuthResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         domain.User `json:"user"`
}

func (s *UserService) Login(ctx context.Context, email, password string) (AuthResult, error) {
	if strings.TrimSpace(email) == "" || password == "" {
		return AuthResult{}, domain.ErrInvalidInput
	}
	u, err := s.users.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if errors.Is(err, domain.ErrNotFound) {
		return AuthResult{}, domain.ErrUnauthorized
	}
	if err != nil {
		return AuthResult{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return AuthResult{}, domain.ErrUnauthorized
	}
	access, err := s.accessToken(u)
	if err != nil {
		return AuthResult{}, err
	}
	refresh, err := randomToken()
	if err != nil {
		return AuthResult{}, err
	}
	if err := s.tokens.Create(ctx, domain.RefreshToken{ID: uuid.New(), UserID: u.ID, TokenHash: hashToken(refresh), ExpiresAt: time.Now().Add(s.refreshTTL)}); err != nil {
		return AuthResult{}, err
	}
	u.PasswordHash = ""
	return AuthResult{AccessToken: access, RefreshToken: refresh, ExpiresIn: int64(s.accessTTL.Seconds()), User: u}, nil
}
func (s *UserService) Refresh(ctx context.Context, refresh string) (AuthResult, error) {
	if refresh == "" {
		return AuthResult{}, domain.ErrInvalidInput
	}
	t, err := s.tokens.GetActiveByHash(ctx, hashToken(refresh))
	if errors.Is(err, domain.ErrNotFound) {
		return AuthResult{}, domain.ErrUnauthorized
	}
	if err != nil {
		return AuthResult{}, err
	}
	u, err := s.users.GetByID(ctx, t.UserID)
	if err != nil {
		return AuthResult{}, err
	}
	access, err := s.accessToken(u)
	if err != nil {
		return AuthResult{}, err
	}
	u.PasswordHash = ""
	return AuthResult{AccessToken: access, ExpiresIn: int64(s.accessTTL.Seconds()), User: u}, nil
}
func (s *UserService) Logout(ctx context.Context, refresh string) error {
	t, err := s.tokens.GetActiveByHash(ctx, hashToken(refresh))
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.tokens.Revoke(ctx, t.ID)
}
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	u, err := s.users.GetByID(ctx, id)
	u.PasswordHash = ""
	return u, err
}
func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, name, phone string) (domain.User, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(phone) == "" {
		return domain.User{}, domain.ErrInvalidInput
	}
	u, err := s.users.UpdateProfile(ctx, id, strings.TrimSpace(name), strings.TrimSpace(phone))
	u.PasswordHash = ""
	return u, err
}
func (s *UserService) ParseAccessToken(raw string) (uuid.UUID, domain.Role, error) {
	token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	role, _ := claims["role"].(string)
	return id, domain.Role(role), nil
}
func (s *UserService) accessToken(u domain.User) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.ID.String(), "role": string(u.Role), "iat": now.Unix(), "exp": now.Add(s.accessTTL).Unix()}).SignedString([]byte(s.secret))
}
func validateUserInput(name, email, phone, password string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(phone) == "" || password == "" || len(password) < 8 {
		return domain.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(email); err != nil || !strings.Contains(email, "@") {
		return domain.ErrInvalidInput
	}
	return nil
}
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func hashToken(v string) string {
	sum := sha256.Sum256([]byte(v))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
