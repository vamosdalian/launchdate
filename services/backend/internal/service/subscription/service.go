package subscription

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const subscriptionCollectionName = "email_subscriptions"

var (
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrInvalidToken       = errors.New("invalid unsubscribe token")
	ErrEmailNotConfigured = errors.New("email service is not configured")
)

type EmailSender interface {
	SendSubscriptionStatusEmail(ctx context.Context, to string, status models.SubscriptionStatus, unsubscribeURL string) error
}

type Service struct {
	mc         *db.MongoDB
	sn         *snowflake.Node
	logger     *logrus.Logger
	sender     EmailSender
	webBaseURL string
}

func NewService(mc *db.MongoDB, logger *logrus.Logger, sender EmailSender, webBaseURL string) *Service {
	node, _ := snowflake.NewNode(0)
	return &Service{
		mc:         mc,
		sn:         node,
		logger:     logger,
		sender:     sender,
		webBaseURL: strings.TrimRight(webBaseURL, "/"),
	}
}

func (s *Service) EnsureIndexes(ctx context.Context) error {
	collection := s.mc.Collection(subscriptionCollectionName)
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
		{
			Keys:    bson.D{{Key: "unsubscribe_token", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("unsubscribe_token_unique"),
		},
	})
	return err
}

func NormalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", ErrInvalidEmail
	}

	parsed, err := mail.ParseAddress(normalized)
	if err != nil || !strings.EqualFold(parsed.Address, normalized) {
		return "", ErrInvalidEmail
	}

	return normalized, nil
}

func (s *Service) Subscribe(ctx context.Context, email string) error {
	if s.sender == nil {
		return ErrEmailNotConfigured
	}

	normalizedEmail, err := NormalizeEmail(email)
	if err != nil {
		return err
	}

	collection := s.mc.Collection(subscriptionCollectionName)
	now := time.Now().UTC()

	var subscription models.EmailSubscription
	err = collection.FindOne(ctx, bson.M{"email": normalizedEmail}).Decode(&subscription)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("find subscription by email: %w", err)
	}

	if err == mongo.ErrNoDocuments {
		token, tokenErr := generateToken()
		if tokenErr != nil {
			return tokenErr
		}

		subscription = models.EmailSubscription{
			ID:                  s.sn.Generate().Int64(),
			Email:               normalizedEmail,
			Status:              models.SubscriptionStatusSubscribed,
			UnsubscribeToken:    token,
			SubscribedAt:        now,
			LastStatusChangedAt: now,
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		if _, err := collection.InsertOne(ctx, subscription); err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}
	} else if subscription.Status == models.SubscriptionStatusUnsubscribed {
		subscription.Status = models.SubscriptionStatusSubscribed
		subscription.SubscribedAt = now
		subscription.UnsubscribedAt = nil
		subscription.LastStatusChangedAt = now
		subscription.UpdatedAt = now

		_, err = collection.UpdateOne(ctx, bson.M{"email": normalizedEmail}, bson.M{
			"$set": bson.M{
				"status":                 subscription.Status,
				"subscribed_at":          subscription.SubscribedAt,
				"unsubscribed_at":        nil,
				"last_status_changed_at": subscription.LastStatusChangedAt,
				"updated_at":             subscription.UpdatedAt,
			},
		})
		if err != nil {
			return fmt.Errorf("restore subscription: %w", err)
		}
	}

	unsubscribeURL := s.buildUnsubscribeURL(subscription.UnsubscribeToken)
	if err := s.sender.SendSubscriptionStatusEmail(ctx, normalizedEmail, models.SubscriptionStatusSubscribed, unsubscribeURL); err != nil {
		s.logger.WithError(err).WithField("email", normalizedEmail).Error("failed to send subscription status email")
		return fmt.Errorf("send subscription email: %w", err)
	}

	return nil
}

func (s *Service) UnsubscribeByToken(ctx context.Context, token string) (models.SubscriptionStatus, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", ErrInvalidToken
	}

	collection := s.mc.Collection(subscriptionCollectionName)
	now := time.Now().UTC()

	var subscription models.EmailSubscription
	err := collection.FindOne(ctx, bson.M{"unsubscribe_token": token}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", ErrInvalidToken
		}
		return "", fmt.Errorf("find subscription by token: %w", err)
	}

	if subscription.Status == models.SubscriptionStatusUnsubscribed {
		return models.SubscriptionStatusUnsubscribed, nil
	}

	_, err = collection.UpdateOne(ctx, bson.M{"unsubscribe_token": token}, bson.M{
		"$set": bson.M{
			"status":                 models.SubscriptionStatusUnsubscribed,
			"unsubscribed_at":        now,
			"last_status_changed_at": now,
			"updated_at":             now,
		},
	})
	if err != nil {
		return "", fmt.Errorf("unsubscribe subscription: %w", err)
	}

	return models.SubscriptionStatusUnsubscribed, nil
}

func (s *Service) RenderUnsubscribeResultPage(status models.SubscriptionStatus, invalidToken bool) string {
	title := "Subscription Update"
	heading := "Subscription updated"
	description := "Your LaunchDate email preferences have been updated."

	switch {
	case invalidToken:
		title = "Invalid Unsubscribe Link"
		heading = "This unsubscribe link is not valid"
		description = "The link may have expired or been copied incorrectly. Return to LaunchDate and subscribe again if needed."
	case status == models.SubscriptionStatusUnsubscribed:
		title = "Unsubscribed"
		heading = "You are unsubscribed"
		description = "Your email address has been removed from LaunchDate subscription updates. No further action is required."
	}

	homeURL := s.webBaseURL
	if homeURL == "" {
		homeURL = "https://launch-date.com"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>%s</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #07111d;
      --panel: #0f1724;
      --panel-border: rgba(148, 163, 184, 0.18);
      --text: #f8fafc;
      --muted: #94a3b8;
      --accent: #2563eb;
      --accent-hover: #1d4ed8;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      font-family: Arial, sans-serif;
      background:
        radial-gradient(circle at top, rgba(37, 99, 235, 0.22), transparent 32%%),
        linear-gradient(180deg, #020617 0%%, var(--bg) 100%%);
      color: var(--text);
    }
    .shell {
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }
    header, footer {
      padding: 24px;
      border-color: var(--panel-border);
      background: rgba(2, 6, 23, 0.82);
      backdrop-filter: blur(18px);
    }
    header { border-bottom: 1px solid var(--panel-border); }
    footer { border-top: 1px solid var(--panel-border); }
    .brand {
      max-width: 1120px;
      margin: 0 auto;
      font-size: 13px;
      font-weight: 700;
      letter-spacing: 0.18em;
      text-transform: uppercase;
      color: #cbd5e1;
    }
    main {
      flex: 1;
      display: grid;
      place-items: center;
      padding: 40px 24px;
    }
    .card {
      width: 100%%;
      max-width: 720px;
      padding: 40px;
      border: 1px solid var(--panel-border);
      border-radius: 24px;
      background: rgba(15, 23, 36, 0.88);
      box-shadow: 0 30px 80px rgba(2, 6, 23, 0.4);
    }
    h1 {
      margin: 0 0 16px;
      font-size: clamp(32px, 6vw, 52px);
      line-height: 1;
    }
    p {
      margin: 0;
      color: var(--muted);
      font-size: 18px;
      line-height: 1.7;
    }
    .actions { margin-top: 28px; }
    .button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 48px;
      padding: 0 20px;
      border-radius: 999px;
      background: var(--accent);
      color: #fff;
      text-decoration: none;
      font-weight: 700;
    }
    .button:hover { background: var(--accent-hover); }
    .footer-inner {
      max-width: 1120px;
      margin: 0 auto;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      color: var(--muted);
      font-size: 14px;
    }
    .footer-link {
      color: #cbd5e1;
      text-decoration: none;
    }
  </style>
</head>
<body>
  <div class="shell">
    <header>
      <div class="brand">LaunchDate</div>
    </header>
    <main>
      <section class="card">
        <h1>%s</h1>
        <p>%s</p>
        <div class="actions">
          <a class="button" href="%s">Return to LaunchDate</a>
        </div>
      </section>
    </main>
    <footer>
      <div class="footer-inner">
        <span>&copy; 2026 LaunchDate. Mission updates and launch tracking.</span>
        <a class="footer-link" href="%s">launch-date.com</a>
      </div>
    </footer>
  </div>
</body>
</html>`,
		html.EscapeString(title),
		html.EscapeString(heading),
		html.EscapeString(description),
		html.EscapeString(homeURL),
		html.EscapeString(homeURL),
	)
}

func (s *Service) buildUnsubscribeURL(token string) string {
	baseURL := s.webBaseURL
	if baseURL == "" {
		baseURL = "https://launch-date.com"
	}
	return fmt.Sprintf("%s/api/v1/subscriptions/unsubscribe?token=%s", baseURL, url.QueryEscape(token))
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate unsubscribe token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
