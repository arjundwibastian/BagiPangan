package handler

import (
	"context"
	"log"

	"github.com/google/uuid"
	userv1 "github.com/zero-hunger/contracts/gen/user/v1"
	"github.com/zero-hunger/request-service/internal/domain"
)

type UserClient struct {
	client userv1.UserServiceClient
	log    *log.Logger
}

func NewUserClient(client userv1.UserServiceClient, logger *log.Logger) *UserClient {
	if logger == nil {
		logger = log.Default()
	}
	return &UserClient{client: client, log: logger}
}

func (c *UserClient) ValidateRecipient(ctx context.Context, id uuid.UUID) error {
	c.log.Printf("event=user_service_get_user_start user_id=%s", id)
	response, err := c.client.GetUser(ctx, &userv1.GetUserRequest{UserId: id.String()})
	if err != nil {
		c.log.Printf("event=user_service_get_user_failed user_id=%s error=%v", id, err)
		return err
	}
	if response.GetRole() != userv1.UserRole_USER_ROLE_RECIPIENT {
		c.log.Printf("event=user_service_recipient_rejected user_id=%s role=%s", id, response.GetRole())
		return domain.ErrForbidden
	}
	c.log.Printf("event=user_service_recipient_validated user_id=%s", id)
	return nil
}
