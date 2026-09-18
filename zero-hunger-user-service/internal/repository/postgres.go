package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zero-hunger/user-service/internal/domain"
)

type UserRepository struct{ db *pgxpool.Pool }

func NewUserRepository(db *pgxpool.Pool) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	err := r.db.QueryRow(ctx, `INSERT INTO users (name,email,phone,password_hash,role) VALUES ($1,$2,$3,$4,$5) RETURNING user_id,created_at,updated_at`, user.Name, user.Email, user.Phone, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, `SELECT user_id,name,email,phone,password_hash,role,created_at,updated_at FROM users WHERE user_id=$1`, id))
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, `SELECT user_id,name,email,phone,password_hash,role,created_at,updated_at FROM users WHERE email=$1`, email))
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id uuid.UUID, name, phone string) (domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, `UPDATE users SET name=$2,phone=$3,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1 RETURNING user_id,name,email,phone,password_hash,role,created_at,updated_at`, id, name, phone))
}

type rowScanner interface{ Scan(...any) error }

func scanUser(row rowScanner) (domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	return u, err
}

type RefreshTokenRepository struct{ db *pgxpool.Pool }

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}
func (r *RefreshTokenRepository) Create(ctx context.Context, token domain.RefreshToken) error {
	_, err := r.db.Exec(ctx, `INSERT INTO refresh_tokens(refresh_token_id,user_id,token_hash,expires_at) VALUES($1,$2,$3,$4)`, token.ID, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}
func (r *RefreshTokenRepository) GetActiveByHash(ctx context.Context, hash string) (domain.RefreshToken, error) {
	var t domain.RefreshToken
	err := r.db.QueryRow(ctx, `SELECT refresh_token_id,user_id,token_hash,expires_at,revoked_at,created_at FROM refresh_tokens WHERE token_hash=$1 AND revoked_at IS NULL AND expires_at>CURRENT_TIMESTAMP`, hash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, domain.ErrNotFound
	}
	return t, err
}
func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=CURRENT_TIMESTAMP WHERE refresh_token_id=$1`, id)
	return err
}
