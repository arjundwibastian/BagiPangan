package grpcClient

import (
	"context"

	foodv1 "github.com/zero-hunger/contracts/gen/food/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func NewUserServiceClient(ctx context.Context, address string) (*UserClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &UserClient{
		conn:   conn,
		client: userv1.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*userv1.GetUserResponse, error) {
	return c.client.GetUser(ctx, &userv1.GetUserRequest{UserId: userID})
}

func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type FoodClient struct {
	conn   *grpc.ClientConn
	client foodv1.FoodServiceClient
}

func NewFoodServiceClient(ctx context.Context, address string) (*FoodClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &FoodClient{
		conn:   conn,
		client: foodv1.NewFoodServiceClient(conn),
	}, nil
}

func (c *FoodClient) GetFoodListing(ctx context.Context, foodListingID string) (*foodv1.GetFoodListingResponse, error) {
	return c.client.GetFoodListing(ctx, &foodv1.GetFoodListingRequest{FoodListingId: foodListingID})
}

func (c *FoodClient) ReserveFoodQuantity(ctx context.Context, foodListingID string, claimID string, quantity int32) (*foodv1.ReserveFoodQuantityResponse, error) {
	return c.client.ReserveFoodQuantity(ctx, &foodv1.ReserveFoodQuantityRequest{
		FoodListingId: foodListingID,
		ClaimId:       claimID,
		Quantity:      quantity,
	})
}

func (c *FoodClient) ReleaseFoodQuantity(ctx context.Context, foodListingID string, claimID string, quantity int32) (*foodv1.ReleaseFoodQuantityResponse, error) {
	return c.client.ReleaseFoodQuantity(ctx, &foodv1.ReleaseFoodQuantityRequest{
		FoodListingId: foodListingID,
		ClaimId:       claimID,
		Quantity:      quantity,
	})
}

func (c *FoodClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
