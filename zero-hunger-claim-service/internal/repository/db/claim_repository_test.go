package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zero-hunger/claim-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testPostgres(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("BAGIPANGAN_CLAIM_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set BAGIPANGAN_CLAIM_TEST_DATABASE_URL to the disposable local database")
	}
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid test database configuration")
	}
	if (cfg.Host != "127.0.0.1" && cfg.Host != "localhost" && cfg.Host != "::1") || cfg.Database != "bagipangan_claim_test" || len(cfg.Fallbacks) != 0 {
		t.Fatal("test database must be loopback bagipangan_claim_test without fallback hosts")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		t.Fatal("disposable PostgreSQL unavailable")
	}
	schema := "audit_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := conn.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		conn.Close(ctx)
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()
		if _, err := conn.Exec(cleanup, "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Error(err)
		}
		conn.Close(cleanup)
	})
	cfg.RuntimeParams["search_path"] = schema
	sqlDB := stdlib.OpenDB(*cfg)
	t.Cleanup(func() { sqlDB.Close() })
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../migrations/001_init.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := strings.ReplaceAll(string(migration), "WHERE typname =", "WHERE typnamespace = current_schema()::regnamespace AND typname =")
	if err := db.Exec(sql).Error; err != nil {
		t.Fatal(err)
	}
	return db
}

func TestClaimRepositoryPostgres(t *testing.T) {
	db := testPostgres(t)
	repo := NewClaimRepository(db)
	ctx := context.Background()
	sequence := 0
	newClaim := func() *domain.FoodClaim {
		t.Helper()
		sequence++
		claim, err := repo.Create(ctx, &domain.FoodClaim{ClaimID: uuid.New(), RequestID: uuid.New(), FoodListingID: uuid.New(),
			ClaimCode: fmt.Sprintf("%06d", sequence), ClaimedQuantity: 1, Status: domain.StatusWaitingForPickup, CodeExpiresAt: time.Now().Add(time.Hour)})
		if err != nil {
			t.Fatal(err)
		}
		return claim
	}
	t.Run("pickup retry preserves timestamp even after code expires", func(t *testing.T) {
		claim := newClaim()
		first, err := repo.VerifyPickupCode(ctx, claim.ClaimID, claim.ClaimCode)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Model(claim).Update("code_expires_at", time.Now().Add(-time.Minute)).Error; err != nil {
			t.Fatal(err)
		}
		retried, err := repo.VerifyPickupCode(ctx, claim.ClaimID, claim.ClaimCode)
		if err != nil || retried.ClaimedAt == nil || !retried.ClaimedAt.Truncate(time.Microsecond).Equal(first.ClaimedAt.Truncate(time.Microsecond)) {
			t.Fatalf("retry=%+v error=%v", retried, err)
		}
		if _, err := repo.VerifyPickupCode(ctx, claim.ClaimID, "999999"); !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("wrong code: %v", err)
		}
	})
	t.Run("expiry prevents initial pickup", func(t *testing.T) {
		claim := newClaim()
		if err := db.Model(claim).Update("code_expires_at", time.Now().Add(-time.Minute)).Error; err != nil {
			t.Fatal(err)
		}
		if _, err := repo.VerifyPickupCode(ctx, claim.ClaimID, claim.ClaimCode); !errors.Is(err, domain.ErrClaimExpired) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("pickup and cancellation cannot both transition", func(t *testing.T) {
		claim := newClaim()
		start := make(chan struct{})
		results := make(chan error, 2)
		go func() { <-start; _, err := repo.VerifyPickupCode(ctx, claim.ClaimID, claim.ClaimCode); results <- err }()
		go func() { <-start; _, err := repo.Cancel(ctx, claim.ClaimID); results <- err }()
		close(start)
		wins := 0
		for range 2 {
			err := <-results
			if err == nil {
				wins++
			} else if !errors.Is(err, domain.ErrClaimAlreadyVerified) {
				t.Fatal(err)
			}
		}
		if wins != 1 {
			t.Fatalf("expected exactly one transition, got %d", wins)
		}
	})
	t.Run("cancelled claim cannot be picked up or cancelled again", func(t *testing.T) {
		claim := newClaim()
		if _, err := repo.Cancel(ctx, claim.ClaimID); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.Cancel(ctx, claim.ClaimID); !errors.Is(err, domain.ErrClaimAlreadyVerified) {
			t.Fatalf("repeat cancel: %v", err)
		}
		if _, err := repo.VerifyPickupCode(ctx, claim.ClaimID, claim.ClaimCode); !errors.Is(err, domain.ErrClaimAlreadyVerified) {
			t.Fatalf("cancelled pickup: %v", err)
		}
	})
}
