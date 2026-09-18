package grpcClient

import (
	"context"

	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type requestClient struct {
	conn   *grpc.ClientConn
	client requestv1.RequestServiceClient
}

func NewRequestClient(ctx context.Context, address string) (*requestClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &requestClient{
		conn:   conn,
		client: requestv1.NewRequestServiceClient(conn),
	}, nil
}

func (c *requestClient) GetRequest(ctx context.Context, requestID string) (*requestv1.GetRequestResponse, error) {
	return c.client.GetRequest(ctx, &requestv1.GetRequestRequest{RequestId: requestID})
}
