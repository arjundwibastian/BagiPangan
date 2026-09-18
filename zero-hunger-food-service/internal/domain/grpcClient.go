package domain

import (
	"context"

	requestv1 "github.com/zero-hunger/contracts/gen/request/v1"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (*userv1.GetUserResponse, error)
}
type RequestService interface {
	GetRequest(ctx context.Context, requestID string) (*requestv1.GetRequestResponse, error)
}
