package grpcClient

import (
	"context"

	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestClient struct {
	conn   *grpc.ClientConn
	client requestv1.RequestServiceClient
}

func NewRequestServiceClient(
	ctx context.Context,
	address string,
) (*RequestClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &RequestClient{
		conn:   conn,
		client: requestv1.NewRequestServiceClient(conn),
	}, nil
}

func (c *RequestClient) GetRequest(
	ctx context.Context,
	requestID string,
) (*requestv1.GetRequestResponse, error) {
	return c.client.GetRequest(
		ctx,
		&requestv1.GetRequestRequest{
			RequestId: requestID,
		},
	)
}

func (c *RequestClient) MarkRequestClaimed(
	ctx context.Context,
	requestID string,
	claimID string,
) (*requestv1.MarkRequestClaimedResponse, error) {
	return c.client.MarkRequestClaimed(
		ctx,
		&requestv1.MarkRequestClaimedRequest{
			RequestId: requestID,
			ClaimId:   claimID,
		},
	)
}

func (c *RequestClient) MarkRequestSearching(
	ctx context.Context,
	requestID string,
	claimID string,
) (*requestv1.MarkRequestSearchingResponse, error) {
	return c.client.MarkRequestSearching(
		ctx,
		&requestv1.MarkRequestSearchingRequest{
			RequestId: requestID,
			ClaimId:   claimID,
		},
	)
}

func (c *RequestClient) MarkRequestCompleted(
	ctx context.Context,
	requestID string,
	claimID string,
) (*requestv1.MarkRequestCompletedResponse, error) {
	return c.client.MarkRequestCompleted(
		ctx,
		&requestv1.MarkRequestCompletedRequest{
			RequestId: requestID,
			ClaimId:   claimID,
		},
	)
}

func (c *RequestClient) Close() error {
	return c.conn.Close()
}
