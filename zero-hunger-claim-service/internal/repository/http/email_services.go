package http

import (
	"fmt"
	"html"

	"github.com/resend/resend-go/v2"
	"github.com/zero-hunger/claim-service/internal/domain"
)

type SendGridEmailService struct {
	apiKey string
}

func NewSendGridEmailService(apiKey string) domain.EmailService {
	return &SendGridEmailService{apiKey: apiKey}
}
func (s *SendGridEmailService) SendClaimNotification(claim *domain.ClaimNotificationData) error {
	client := resend.NewClient(s.apiKey)

	recipientName := html.EscapeString(claim.RecipientName)
	foodName := html.EscapeString(claim.FoodName) 
	claimCode := html.EscapeString(claim.ClaimCode)
	status := html.EscapeString(string(claim.Status))

	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body style="margin:0; padding:0; background-color:#f4f7f2; font-family:Arial,sans-serif; color:#263326;">
		<div style="max-width:620px; margin:30px auto; background:#ffffff; border-radius:16px; overflow:hidden; box-shadow:0 5px 20px rgba(0,0,0,0.08);">

			<div style="background:linear-gradient(135deg,#2e7d32,#66bb6a); padding:30px; text-align:center; color:white;">
				<div style="font-size:42px;">🥗</div>
				<h1 style="margin:10px 0 5px; font-size:28px;">Food Claim Confirmed!</h1>
				<p style="margin:0; font-size:15px;">Good food is on its way to someone who needs it.</p>
			</div>

			<div style="padding:30px;">
				<p style="font-size:16px;">Hello <strong>%s</strong>,</p>

				<p style="line-height:1.6;">
					Your food claim has been successfully created. Please keep your claim code safe and show it when collecting your food.
				</p>

				<div style="background:#f1f8e9; border:1px solid #c5e1a5; border-radius:12px; padding:22px; text-align:center; margin:25px 0;">
					<p style="margin:0 0 8px; color:#558b2f; font-size:13px; text-transform:uppercase; letter-spacing:1px;">
						Your Claim Code
					</p>

					<div style="font-size:34px; font-weight:bold; letter-spacing:6px; color:#33691e;">
						%s
					</div>

					<p style="margin:12px 0 0; font-size:13px; color:#66705f;">
						Expires on %s
					</p>
				</div>

				<h3 style="color:#33691e; margin-bottom:15px;">Claim Details</h3>

				<table style="width:100%%; border-collapse:collapse; font-size:15px;">
					<tr>
						<td style="padding:12px 0; border-bottom:1px solid #eeeeee; color:#777;">Food</td>
						<td style="padding:12px 0; border-bottom:1px solid #eeeeee; text-align:right; font-weight:bold;">%s</td>
					</tr>
					<tr>
						<td style="padding:12px 0; border-bottom:1px solid #eeeeee; color:#777;">Quantity</td>
						<td style="padding:12px 0; border-bottom:1px solid #eeeeee; text-align:right; font-weight:bold;">%d portion(s)</td>
					</tr>
					<tr>
						<td style="padding:12px 0; color:#777;">Status</td>
						<td style="padding:12px 0; text-align:right; font-weight:bold; color:#2e7d32;">%s</td>
					</tr>
				</table>

				<div style="background:#fff8e1; border-left:4px solid #f9a825; padding:14px; margin-top:25px; font-size:14px; line-height:1.5;">
					<strong>Important:</strong> Please collect your food before the claim code expires.
					Bring this code when picking up your food.
				</div>

				<p style="margin-top:30px; line-height:1.6;">
					Thank you for helping reduce food waste and making a positive impact in the community. 💚
				</p>

				<p style="margin-bottom:0;">
					Best regards,<br>
					<strong>Zero Hunger Team</strong>
				</p>
			</div>

			<div style="background:#f8faf7; padding:18px; text-align:center; color:#899486; font-size:12px;">
				This is an automated message. Please do not reply directly to this email.
			</div>
		</div>
	</body>
	</html>
	`,
		recipientName,
		claimCode,
		claim.CodeExpiresAt.Format("02 Jan 2006 at 15:04"),
		foodName,
		claim.Quantity,
		status,
	)

params := &resend.SendEmailRequest{
	From:    "Zero Hunger <admin@resend.dev>",
	To:      []string{claim.RecipientEmail},
	Subject: "🥗 Your Food Claim Is Confirmed",
	Html:    htmlContent,
}
    // Send it!
    _, err := client.Emails.Send(params)
    if err != nil {
        return err
    }
    
	return nil
}


func (s *SendGridEmailService) SendCancelNotification(claim *domain.ClaimNotificationData) error  {
	client := resend.NewClient(s.apiKey)

	recipientName := html.EscapeString(claim.RecipientName)
	recipientEmail := html.EscapeString(claim.RecipientEmail)
	claimCode := html.EscapeString(claim.ClaimCode)
	status := html.EscapeString(claim.Status)

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body style="margin:0; padding:0; background:#f4f7f2; font-family:Arial,sans-serif; color:#263326;">
			<div style="max-width:620px; margin:30px auto; background:#ffffff; border-radius:16px; overflow:hidden; box-shadow:0 5px 20px rgba(0,0,0,0.08);">

				<div style="background:linear-gradient(135deg,#ef6c00,#ffb74d); padding:30px; text-align:center; color:white;">
					<div style="font-size:42px;">📦</div>
					<h1 style="margin:10px 0 5px; font-size:28px;">Claim Cancelled</h1>
					<p style="margin:0; font-size:15px;">Your claim has been cancelled successfully.</p>
				</div>

				<div style="padding:30px;">
					<p style="font-size:16px;">Hello <strong>%s</strong>,</p>

					<p style="line-height:1.6;">
						Your food claim cancellation has been processed. The reserved quantity has been released,
						and your request is now available for searching again.
					</p>

					<div style="background:#fff8e1; border:1px solid #ffcc80; border-radius:12px; padding:22px; margin:25px 0;">
						<h3 style="margin-top:0; color:#e65100;">Claim Information</h3>

						<table style="width:100%%; border-collapse:collapse; font-size:15px;">
							<tr>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; color:#777;">
									Recipient
								</td>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; text-align:right; font-weight:bold;">
									%s
								</td>
							</tr>

							<tr>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; color:#777;">
									Email
								</td>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; text-align:right;">
									%s
								</td>
							</tr>

							<tr>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; color:#777;">
									Claim Code
								</td>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; text-align:right; font-weight:bold; color:#e65100;">
									%s
								</td>
							</tr>

							<tr>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; color:#777;">
									Quantity
								</td>
								<td style="padding:12px 0; border-bottom:1px solid #eeeeee; text-align:right; font-weight:bold;">
									%d portion(s)
								</td>
							</tr>

							<tr>
								<td style="padding:12px 0; color:#777;">
									Status
								</td>
								<td style="padding:12px 0; text-align:right; font-weight:bold; color:#e65100;">
									%s
								</td>
							</tr>
						</table>
					</div>

					<p style="line-height:1.6;">
						You may continue searching for another available food request.
					</p>

					<p style="margin-top:28px; line-height:1.6;">
						Thank you for using Zero Hunger and helping reduce food waste. 💚
					</p>

					<p>
						Best regards,<br>
						<strong>Zero Hunger Team</strong>
					</p>
				</div>

				<div style="background:#f8faf7; padding:18px; text-align:center; color:#899486; font-size:12px;">
					This is an automated message. Please do not reply directly to this email.
				</div>
			</div>
		</body>
		</html>
		`,
				recipientName,
				recipientName,
				recipientEmail,
				claimCode,
				claim.Quantity,
				status,
			)

	params := &resend.SendEmailRequest{
		From:    "Zero Hunger <admin@resend.dev>",
		To:      []string{claim.RecipientEmail},
		Subject: "📦 Your Claim Has Been Cancelled",
		Html:    htmlContent,
	}

    // Send it!
    _, err := client.Emails.Send(params)
    if err != nil {
        return err
    }
    
	return nil
}