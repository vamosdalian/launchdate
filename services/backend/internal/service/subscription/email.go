package subscription

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/resend/resend-go/v3"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type ResendEmailSender struct {
	client *resend.Client
	from   string
}

func NewResendEmailSender(apiKey, from string) *ResendEmailSender {
	apiKey = strings.TrimSpace(apiKey)
	from = strings.TrimSpace(from)
	if apiKey == "" || from == "" {
		return nil
	}

	return &ResendEmailSender{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (s *ResendEmailSender) SendSubscriptionStatusEmail(ctx context.Context, to string, status models.SubscriptionStatus, unsubscribeURL string) error {
	if s == nil || s.client == nil {
		return ErrEmailNotConfigured
	}

	subject := "LaunchDate subscription status"
	statusText := "Your subscription status is active."
	bodyText := "You will stay on the LaunchDate subscription list for future mission updates."

	if status == models.SubscriptionStatusUnsubscribed {
		statusText = "Your subscription status is unsubscribed."
		bodyText = "You will no longer receive LaunchDate subscription updates."
	}

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    buildSubscriptionStatusEmailHTML(statusText, bodyText, unsubscribeURL),
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("resend email: %w", err)
	}

	return nil
}

func buildSubscriptionStatusEmailHTML(statusText, bodyText, unsubscribeURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>LaunchDate Subscription Status</title>
</head>
<body style="margin:0;padding:0;background-color:#020617;font-family:Arial,sans-serif;color:#e2e8f0;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#020617;">
    <tr>
      <td align="center">
        <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:720px;">
          <tr>
            <td style="padding:28px 32px;background:#0f172a;border-bottom:1px solid rgba(148,163,184,0.16);">
              <div style="font-size:12px;letter-spacing:0.24em;text-transform:uppercase;font-weight:700;color:#cbd5e1;">LaunchDate</div>
              <div style="margin-top:12px;font-size:28px;line-height:1.2;font-weight:700;color:#f8fafc;">Subscription status update</div>
            </td>
          </tr>
          <tr>
            <td style="padding:40px 32px;background:#020617;">
              <p style="margin:0 0 16px;font-size:20px;line-height:1.5;color:#f8fafc;font-weight:700;">%s</p>
              <p style="margin:0;font-size:16px;line-height:1.8;color:#cbd5e1;">%s</p>
            </td>
          </tr>
          <tr>
            <td style="padding:24px 32px;border-top:1px solid rgba(148,163,184,0.12);background:#0f172a;">
              <p style="margin:0;font-size:13px;line-height:1.7;color:#94a3b8;">
                Manage this subscription anytime:
                <a href="%s" style="color:#93c5fd;text-decoration:underline;">Unsubscribe</a>
              </p>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 32px;border-top:1px solid rgba(148,163,184,0.12);background:#020617;">
              <p style="margin:0;font-size:12px;line-height:1.7;color:#64748b;">LaunchDate mission alerts and launch tracking.</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		html.EscapeString(statusText),
		html.EscapeString(bodyText),
		html.EscapeString(unsubscribeURL),
	)
}
