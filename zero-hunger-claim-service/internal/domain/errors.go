package domain

import "errors"

var (
	ErrNotFound             = errors.New("record not found")
	ErrInvalidInput         = errors.New("invalid input provided")
	ErrFailedToReachService = errors.New("failed to reach destination service")
	ErrServiceError         = errors.New("destination service error")
	ErrFailedDecode         = errors.New("failed to decode service response")
	ErrFailedToGenerateUUID = errors.New("failed to create UUID")
	ErrForbidden            = errors.New("you arent allowed to do this action")
	ErrClaimExpired         = errors.New("claim code has expired")
	ErrClaimAlreadyVerified = errors.New("claim has already been picked up or cancelled")
	ErrInsufficientQuantity = errors.New("insufficient food quantity available")
)
