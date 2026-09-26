package handler

import (
	"context"

	"github.com/google/uuid"
	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	"github.com/zero-hunger/request-service/internal/domain"
	"github.com/zero-hunger/request-service/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RequestRPCServer struct {
	requestv1.UnimplementedRequestServiceServer
	Requests *usecase.RequestService
}

func (s *RequestRPCServer) CreateRequest(ctx context.Context, in *requestv1.CreateRequestRequest) (*requestv1.CreateRequestResponse, error) {
	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	out, err := s.Requests.Create(ctx, usecase.CreateInput{
		UserID:    userID,
		Latitude:  in.GetLatitude(),
		Longitude: in.GetLongitude(),
		RadiusKM:  in.GetRadiusKm()},
	)
	if err != nil {
		return nil, err
	}
	return &requestv1.CreateRequestResponse{Request: toProtoRequest(out.Request)}, nil
}

func (s *RequestRPCServer) GetRequest(ctx context.Context, in *requestv1.GetRequestRequest) (*requestv1.GetRequestResponse, error) {
	id, err := uuid.Parse(in.GetRequestId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	r, err := s.Requests.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &requestv1.GetRequestResponse{Request: toProtoRequest(r)}, nil
}

func (s *RequestRPCServer) ListRequestsByUser(ctx context.Context, in *requestv1.ListRequestsByUserRequest) (*requestv1.ListRequestsByUserResponse, error) {
	id, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	requests, err := s.Requests.ListByUser(ctx, id)
	if err != nil {
		return nil, err
	}
	out := &requestv1.ListRequestsByUserResponse{Requests: make([]*requestv1.FoodRequest, 0, len(requests))}
	for _, r := range requests {
		out.Requests = append(out.Requests, toProtoRequest(r))
	}
	return out, nil
}

func (s *RequestRPCServer) CancelRequest(ctx context.Context, in *requestv1.CancelRequestRequest) (*requestv1.CancelRequestResponse, error) {
	id, err := uuid.Parse(in.GetRequestId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	r, err := s.Requests.Cancel(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return &requestv1.CancelRequestResponse{Request: toProtoRequest(r)}, nil
}

func (s *RequestRPCServer) MarkRequestClaimed(ctx context.Context, in *requestv1.MarkRequestClaimedRequest) (*requestv1.MarkRequestClaimedResponse, error) {
	id, err := uuid.Parse(in.GetRequestId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	claimID, err := uuid.Parse(in.GetClaimId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	r, err := s.Requests.MarkClaimed(ctx, id, claimID)
	if err != nil {
		return nil, err
	}
	return &requestv1.MarkRequestClaimedResponse{Request: toProtoRequest(r)}, nil
}

func (s *RequestRPCServer) MarkRequestCompleted(ctx context.Context, in *requestv1.MarkRequestCompletedRequest) (*requestv1.MarkRequestCompletedResponse, error) {
	id, err := uuid.Parse(in.GetRequestId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	claimID, err := uuid.Parse(in.GetClaimId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	r, err := s.Requests.MarkCompleted(ctx, id, claimID)
	if err != nil {
		return nil, err
	}
	return &requestv1.MarkRequestCompletedResponse{Request: toProtoRequest(r)}, nil
}

func (s *RequestRPCServer) MarkRequestSearching(ctx context.Context, in *requestv1.MarkRequestSearchingRequest) (*requestv1.MarkRequestSearchingResponse, error) {
	id, err := uuid.Parse(in.GetRequestId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	claimID, err := uuid.Parse(in.GetClaimId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	r, err := s.Requests.MarkSearching(ctx, id, claimID)
	if err != nil {
		return nil, err
	}
	return &requestv1.MarkRequestSearchingResponse{Request: toProtoRequest(r)}, nil
}

func toProtoRequest(r domain.FoodRequest) *requestv1.FoodRequest {
	out := &requestv1.FoodRequest{
		RequestId: r.ID.String(), UserId: r.UserID.String(),
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		RadiusKm:  r.RadiusKM,
		Status:    statusToProto(r.Status),
		CreatedAt: timestamppb.New(r.CreatedAt),
	}
	return out
}

func statusToProto(status domain.RequestStatus) requestv1.RequestStatus {
	switch status {
	case domain.StatusSearching:
		return requestv1.RequestStatus_REQUEST_STATUS_SEARCHING
	case domain.StatusClaimed:
		return requestv1.RequestStatus_REQUEST_STATUS_CLAIMED
	case domain.StatusCompleted:
		return requestv1.RequestStatus_REQUEST_STATUS_COMPLETED
	case domain.StatusCancelled:
		return requestv1.RequestStatus_REQUEST_STATUS_CANCELLED
	case domain.StatusExpired:
		return requestv1.RequestStatus_REQUEST_STATUS_EXPIRED
	default:
		return requestv1.RequestStatus_REQUEST_STATUS_UNSPECIFIED
	}
}

var _ requestv1.RequestServiceServer = (*RequestRPCServer)(nil)
