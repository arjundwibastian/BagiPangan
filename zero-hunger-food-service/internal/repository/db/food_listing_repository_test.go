package db

import (
	"context"
	"errors"
	"math"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zero-hunger/food-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Opt-in only: refuses non-loopback hosts and any database except our fixture DB.
// A random schema prevents concurrent test runs from sharing or deleting rows.
func testPostgres(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("BAGIPANGAN_FOOD_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set BAGIPANGAN_FOOD_TEST_DATABASE_URL to the disposable local database")
	}
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid test database configuration")
	}
	if (cfg.Host != "127.0.0.1" && cfg.Host != "localhost" && cfg.Host != "::1") || cfg.Database != "bagipangan_food_test" || len(cfg.Fallbacks) != 0 {
		t.Fatal("test database must be loopback bagipangan_food_test without fallback hosts")
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
		_, err := conn.Exec(cleanup, "DROP SCHEMA "+quoted+" CASCADE")
		if err != nil {
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
	// Use service migration, scoped to this schema. Type existence checks must
	// also be schema scoped so independently running fixtures do not collide.
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

func TestFoodRepositoryPostgres(t *testing.T) {
	db := testPostgres(t)
	repo := NewFoodListingRepo(db)
	ctx := context.Background()
	listing := func(quantity int32, status domain.ListingStatus, expired bool) *domain.FoodListing {
		t.Helper()
		until := time.Now().Add(time.Hour)
		if expired {
			until = time.Now().Add(-time.Minute)
		}
		food, err := repo.InsertNewFoodListing(ctx, &domain.FoodListing{
			FoodListingID: uuid.New(), UserID: uuid.New(), Title: "Isolated fixture", FoodType: domain.FoodCooked,
			Quantity: quantity, Status: status, AvailableFrom: time.Now().Add(-time.Hour), AvailableUntil: until,
		})
		if err != nil {
			t.Fatal(err)
		}
		return food
	}
	t.Run("only one concurrent last-portion reservation wins", func(t *testing.T) {
		food := listing(1, domain.StatusAvailable, false)
		start := make(chan struct{})
		results := make(chan error, 12)
		var wg sync.WaitGroup
		for range 12 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				_, err := repo.ReserveFoodQuantity(ctx, food.FoodListingID, 1)
				results <- err
			}()
		}
		close(start)
		wg.Wait()
		close(results)
		wins := 0
		for err := range results {
			if err == nil {
				wins++
			} else if !errors.Is(err, domain.ErrListingUnavailable) && !errors.Is(err, domain.ErrInsufficientQuantity) {
				t.Fatal(err)
			}
		}
		current, err := repo.GetFoodListingByID(ctx, food.FoodListingID)
		if err != nil || wins != 1 || current.Quantity != 0 {
			t.Fatalf("wins=%d listing=%+v error=%v", wins, current, err)
		}
	})
	t.Run("release preserves donor cancellation", func(t *testing.T) {
		food := listing(3, domain.StatusAvailable, false)
		if _, err := repo.CancelFoodListing(ctx, food.FoodListingID); err != nil {
			t.Fatal(err)
		}
		current, err := repo.ReleaseFoodQuantity(ctx, food.FoodListingID, 2)
		if err != nil || current.Status != domain.StatusExpired || current.Quantity != 5 {
			t.Fatalf("listing=%+v error=%v", current, err)
		}
	})
	t.Run("fully reserved food still expires", func(t *testing.T) {
		food := listing(0, domain.StatusClaimed, true)
		if err := repo.UpdateExpiredAndClaimedListings(ctx); err != nil {
			t.Fatal(err)
		}
		current, err := repo.GetFoodListingByID(ctx, food.FoodListingID)
		if err != nil || current.Status != domain.StatusExpired {
			t.Fatalf("listing=%+v error=%v", current, err)
		}
	})
	t.Run("stale server edit cannot overwrite reservation", func(t *testing.T) {
		food := listing(5, domain.StatusAvailable, false)
		if _, err := repo.ReserveFoodQuantity(ctx, food.FoodListingID, 1); err != nil {
			t.Fatal(err)
		}
		quantity := int32(5)
		_, err := repo.UpdateFoodListing(ctx, food.FoodListingID, domain.FoodListingUpdate{Quantity: &quantity, ExpectedUpdatedAt: food.UpdatedAt})
		if !errors.Is(err, domain.ErrListingUnavailable) {
			t.Fatalf("expected conflict, got %v", err)
		}
		current, err := repo.GetFoodListingByID(ctx, food.FoodListingID)
		if err != nil || current.Quantity != 4 {
			t.Fatalf("listing=%+v error=%v", current, err)
		}
	})
	t.Run("cancel cannot override a fully reserved listing", func(t *testing.T) {
		food := listing(0, domain.StatusClaimed, false)
		if _, err := repo.CancelFoodListing(ctx, food.FoodListingID); !errors.Is(err, domain.ErrListingUnavailable) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("release rejects integer overflow", func(t *testing.T) {
		food := listing(math.MaxInt32, domain.StatusAvailable, false)
		if _, err := repo.ReleaseFoodQuantity(ctx, food.FoodListingID, 1); !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("cancelled query context is propagated", func(t *testing.T) {
		cancelled, cancel := context.WithCancel(ctx)
		cancel()
		if _, err := repo.GetActiveFoodListing(cancelled, string(domain.StatusAvailable)); !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v", err)
		}
	})
}
