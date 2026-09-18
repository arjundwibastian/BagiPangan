package domain

type EmailService interface {
	SendClaimNotification(claim *ClaimNotificationData) error
	SendCancelNotification(claim *ClaimNotificationData) error
}