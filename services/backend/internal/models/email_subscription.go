package models

import "time"

type SubscriptionStatus string

const (
	SubscriptionStatusSubscribed   SubscriptionStatus = "subscribed"
	SubscriptionStatusUnsubscribed SubscriptionStatus = "unsubscribed"
)

type EmailSubscription struct {
	ID                  int64              `json:"id" bson:"id"`
	Email               string             `json:"email" bson:"email"`
	Status              SubscriptionStatus `json:"status" bson:"status"`
	UnsubscribeToken    string             `json:"unsubscribe_token" bson:"unsubscribe_token"`
	SubscribedAt        time.Time          `json:"subscribed_at" bson:"subscribed_at"`
	UnsubscribedAt      *time.Time         `json:"unsubscribed_at,omitempty" bson:"unsubscribed_at,omitempty"`
	LastStatusChangedAt time.Time          `json:"last_status_changed_at" bson:"last_status_changed_at"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
}

type SubscribeRequest struct {
	Email string `json:"email" binding:"required"`
}
