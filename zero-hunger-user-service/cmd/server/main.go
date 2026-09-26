package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/user-service/internal/config"
	"github.com/zero-hunger/user-service/internal/handler"
	"github.com/zero-hunger/user-service/internal/repository"
	"github.com/zero-hunger/user-service/internal/usecase"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	cfg := config.Load()
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	users := repository.NewUserRepository(db)
	tokens := repository.NewRefreshTokenRepository(db)
	svc := usecase.NewUserService(users, tokens, cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	h := handler.NewHTTPHandler(svc)
	e := echo.New()

	v1 := e.Group("/api/v1")

	// Auth routes
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login", h.Login)
	v1.POST("/auth/refresh", h.Refresh)
	v1.POST("/auth/logout", h.Logout)

	// User routes
	v1.GET("/users/me", h.Me)
	v1.PATCH("/users/me", h.UpdateMe)
	v1.GET("/users/:user_id", h.GetUser)

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, &handler.UserRPCServer{Users: svc})

	go func() {
		if err := e.Start(":" + cfg.ServerPort); err != nil {
			log.Printf("http: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("user service HTTP :%s, gRPC :%s", cfg.ServerPort, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
	_ = fmt.Sprint(time.Now())
}
