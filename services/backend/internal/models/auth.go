package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"` // never expose in JSON
	Role         string             `json:"role" bson:"role"`       // e.g., "admin", "user"
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// RefreshToken represents a refresh token stored in database
type RefreshToken struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	TokenHash  string             `json:"-" bson:"token_hash"` // store hash, not plaintext
	ExpiresAt  time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	RevokedAt  *time.Time         `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	LastUsedAt *time.Time         `json:"last_used_at,omitempty" bson:"last_used_at,omitempty"`
	UserAgent  string             `json:"user_agent" bson:"user_agent"`
	IPAddress  string             `json:"ip_address" bson:"ip_address"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	User        *User  `json:"user,omitempty"`
}

// RefreshResponse represents token refresh response
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}
