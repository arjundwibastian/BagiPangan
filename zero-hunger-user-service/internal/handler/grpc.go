package handler

import (
	"context"

	"github.com/google/uuid"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/user-service/internal/domain"
	"github.com/zero-hunger/user-service/internal/service"
)

type UserRPCServer struct {
	userv1.UnimplementedUserServiceServer
	Users *service.UserService
}

func (s *UserRPCServer) GetUser(ctx context.Context, in *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	id, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	u, err := s.Users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &userv1.GetUserResponse{
		UserId: u.ID.String(),
		Name:   u.Name,
		Email:  u.Email,
		Phone:  u.Phone,
		Role:   userRoleToProto(u.Role),
	}, nil
}

func userRoleToProto(role domain.Role) userv1.UserRole {
	switch role {
	case domain.RoleAdmin:
		return userv1.UserRole_USER_ROLE_ADMIN
	case domain.RoleRecipient:
		return userv1.UserRole_USER_ROLE_RECIPIENT
	case domain.RoleDonor:
		return userv1.UserRole_USER_ROLE_DONOR
	default:
		return userv1.UserRole_USER_ROLE_UNSPECIFIED
	}
}
