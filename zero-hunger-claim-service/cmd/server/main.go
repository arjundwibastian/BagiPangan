package main

import (
	"context"
	"errors"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/zero-hunger/claim-service/docs/swagger"
	"github.com/zero-hunger/claim-service/internal/config"
	"github.com/zero-hunger/claim-service/internal/handler"
	"github.com/zero-hunger/claim-service/internal/infrastructure/database"
	grpcClient "github.com/zero-hunger/claim-service/internal/infrastructure/grpcClient"
	"github.com/zero-hunger/claim-service/internal/middleware"
	"github.com/zero-hunger/claim-service/internal/repository/db"
	"github.com/zero-hunger/claim-service/internal/repository/http"
	"github.com/zero-hunger/claim-service/internal/usecase"
)

// @title Zero Hunger Claim Service API
// @version 1.0
// @description Food claim creation, cancellation, and pickup verification API.
// @host localhost:8084
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
		log.Fatal(err)
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	userClient, err := grpcClient.NewUserServiceClient(
		ctx,
		cfg.GRPC.UserService,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer userClient.Close()

	foodClient, err := grpcClient.NewFoodServiceClient(
		ctx,
		cfg.GRPC.FoodService,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer foodClient.Close()

	requestClient, err := grpcClient.NewRequestServiceClient(
		ctx,
		cfg.GRPC.RequestService,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer requestClient.Close()

	claimRepository := db.NewClaimRepository(dbConn)
	emailService := http.NewSendGridEmailService(cfg.App.ResendAPIKey)

	claimUseCase := usecase.NewClaimUseCase(
		claimRepository,
		userClient,
		requestClient,
		foodClient,
		emailService,
	)

	httpHandler := handler.NewHTTPHandler(claimUseCase)

	e := echo.New()
	e.GET("/healthz", handler.Health)
	e.GET("/readyz", handler.Ready(sqlDB.PingContext))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())

	auth := middleware.Authenticate(cfg.JWT.Secret)

	e.POST("/api/v1/claims", httpHandler.CreateClaim, auth)

	e.GET("/api/v1/claims/:claim_id", httpHandler.GetClaim, auth)

	e.POST("/api/v1/claims/:claim_id/verify-pickup", httpHandler.VerifyPickupCode, auth)

	e.POST("/api/v1/claims/:claim_id/cancel", httpHandler.CancelClaim, auth)

	go func() {
		if err := e.Start(":" + cfg.App.Port); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()
	log.Printf("claim service listening on :%s", cfg.App.Port)

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-stopCtx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	log.Println("shutdown complete")
}
