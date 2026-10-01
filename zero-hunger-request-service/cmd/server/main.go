package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/request-service/internal/config"
	"github.com/zero-hunger/request-service/internal/handler"
	"github.com/zero-hunger/request-service/internal/repository"
	"github.com/zero-hunger/request-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	userConn, err := grpc.Dial(cfg.UserServiceGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	logger := log.Default()
	repo := repository.NewRequestRepository(db)
	svc := usecase.NewRequestService(repo, handler.NewUserClient(userv1.NewUserServiceClient(userConn), logger), cfg.RequestTTL)
	h := handler.NewHTTPHandler(svc, cfg.JWTSecret)
	e := echo.New()
	e.GET("/healthz", handler.Health)
	e.GET("/readyz", handler.Ready(db.Ping))
	v1 := e.Group("/api/v1")
	v1.POST("/food-requests", h.Create)
	v1.GET("/food-requests/:request_id", h.Get)
	v1.GET("/users/:user_id/food-requests", h.ListByUser)
	v1.PATCH("/food-requests/:request_id/cancel", h.Cancel)

	grpcServer := grpc.NewServer()
	requestv1.RegisterRequestServiceServer(grpcServer, &handler.RequestRPCServer{Requests: svc})
	go func() {
		if err := e.Start(":" + cfg.ServerPort); err != nil {
			log.Printf("http: %v", err)
		}
	}()
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("request service HTTP :%s, gRPC :%s", cfg.ServerPort, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
