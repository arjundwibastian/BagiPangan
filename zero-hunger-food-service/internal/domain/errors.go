package domain

import "errors"

var (
	ErrNotFound             = errors.New("record not found")
	ErrInvalidInput         = errors.New("invalid input provided")
	ErrFailedToReachService = errors.New("Failed to reach destionation service")
	ErrServiceError         = errors.New("destionation service error")
	ErrFailedDecode         = errors.New("failed to decode service response")
	ErrFailedToGenerateUUID = errors.New("failed to create UUID")
	ErrForbidden            = errors.New("you arent allowed to do this action")
	ErrListingUnavailable   = errors.New("food listing is unavailable")
	ErrInsufficientQuantity = errors.New("food listing doesnt have enough requested quantity")
)
