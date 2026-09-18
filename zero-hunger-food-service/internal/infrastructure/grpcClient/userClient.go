package grpcClient

import (
	"context"

	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func NewUserClient(ctx context.Context, address string) (*UserClient, error) {
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
