package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	foodv1 "github.com/zero-hunger/contracts/gen/food/v1"
	_ "github.com/zero-hunger/food-service/docs/swagger"
	"github.com/zero-hunger/food-service/internal/config"
	"github.com/zero-hunger/food-service/internal/domain"
	"github.com/zero-hunger/food-service/internal/handler"
	"github.com/zero-hunger/food-service/internal/infrastructure/database"
	"github.com/zero-hunger/food-service/internal/infrastructure/grpcClient"
	"github.com/zero-hunger/food-service/internal/middleware"
	"github.com/zero-hunger/food-service/internal/repository/db"
	"github.com/zero-hunger/food-service/internal/usecase"
	"github.com/zero-hunger/food-service/internal/validator"
	"google.golang.org/grpc"
)

// @title Zero Hunger Food Service API
// @version 1.0
// @description Food listing management and nearby food search API.
// @host localhost:8082
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token using: Bearer {token}
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbConn, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatal("Error getting database handle:", err)
	}
	jwtSecret := []byte(cfg.JWT.Secret)
	jwtMiddleware := middleware.NewJWTMiddleware(jwtSecret)

	userClient, err := grpcClient.NewUserClient(context.Background(), cfg.GRPC.UserService)
	if err != nil {
		log.Fatal(err)
	}
	requestClient, err := grpcClient.NewRequestClient(context.Background(), cfg.GRPC.RequestService)
	if err != nil {
		log.Fatal(err)
	}

	repo := db.NewFoodListingRepo(dbConn)
	usecase := usecase.NewfoodListingUseCase(repo, userClient, requestClient)
	httpHandler := handler.NewFoodListingHandler(usecase)

	e := echo.New()
	e.GET("/healthz", handler.Health)
	e.GET("/readyz", handler.Ready(sqlDB.PingContext))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())
	e.Validator = validator.NewValidator()

	donor := e.Group("")
	donor.Use(jwtMiddleware.Authenticate, jwtMiddleware.DonorOnly)
	donor.POST("/api/v1/food-listings", httpHandler.CreateFoodListingHandler)
	// e.PUT("/api/v1/orders/:id", orderHandler.UpdateOrderByIDHandler)
	// e.DELETE("/api/v1/orders/:id", orderHandler.DeleteOrderByIDHandler)

	allRole := e.Group("")
	allRole.Use(jwtMiddleware.Authenticate)
	allRole.GET("/api/v1/food-listings", httpHandler.GetActiveFoodListingHandler)
	allRole.GET("/api/v1/food-listings/nearby", httpHandler.GetNearbyFoodListingHandler)
	allRole.GET("/api/v1/food-listings/:food_listing_id", httpHandler.GetFoodListingByIDHandler)

	// ticker for auto updates
	go startListingStatusUpdater(repo)

	grpcServer := grpc.NewServer()
	foodv1.RegisterFoodServiceServer(grpcServer, handler.NewFoodRpcServer(usecase))

	lis, err := net.Listen("tcp", ":"+cfg.App.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}
	go func() {
		log.Printf("gRPC server listening on :%s", cfg.App.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		if err := e.Start(":" + cfg.App.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-stopCtx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	grpcServer.GracefulStop()
	log.Println("shutdown complete")
}




func startListingStatusUpdater(repo domain.FoodListingRepository) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		if err := repo.UpdateExpiredAndClaimedListings(context.Background()); err != nil {
			log.Printf("failed to update listing statuses: %v", err)
		}

		<-ticker.C
	}
}
