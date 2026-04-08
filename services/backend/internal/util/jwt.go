package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(
	secretKey string,
	accessTokenExpireMin int,
	refreshTokenExpireDays int,
	issuer string,
) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		accessTokenDuration:  time.Duration(accessTokenExpireMin) * time.Minute,
		refreshTokenDuration: time.Duration(refreshTokenExpireDays) * 24 * time.Hour,
		issuer:               issuer,
	}
}

// GenerateAccessToken generates a new JWT access token
func (m *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"iss":      m.issuer,
		"iat":      now.Unix(),
		"exp":      now.Add(m.accessTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// GenerateRefreshToken generates a new refresh token (random string)
func (m *JWTManager) GenerateRefreshToken() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"type": "refresh",
		"iat":  now.Unix(),
		"exp":  now.Add(m.refreshTokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// VerifyAccessToken verifies the access token and returns claims
func (m *JWTManager) VerifyAccessToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// verify issuer
	if iss, ok := claims["iss"].(string); !ok || iss != m.issuer {
		return nil, ErrInvalidToken
	}

	// extract claims
	jwtClaims := &models.JWTClaims{
		Issuer: m.issuer,
	}

	if userID, ok := claims["user_id"].(string); ok {
		jwtClaims.UserID = userID
	}
	if username, ok := claims["username"].(string); ok {
		jwtClaims.Username = username
	}
	if role, ok := claims["role"].(string); ok {
		jwtClaims.Role = role
	}
	if iat, ok := claims["iat"].(float64); ok {
		jwtClaims.IssuedAt = int64(iat)
	}
	if exp, ok := claims["exp"].(float64); ok {
		jwtClaims.ExpiresAt = int64(exp)
	}

	return jwtClaims, nil
}

// VerifyRefreshToken verifies the refresh token
func (m *JWTManager) VerifyRefreshToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrExpiredToken
		}
		return ErrInvalidToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	// verify type
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return ErrInvalidToken
	}

	return nil
}

// GetRefreshTokenExpiration returns the refresh token expiration duration
func (m *JWTManager) GetRefreshTokenExpiration() time.Duration {
	return m.refreshTokenDuration
}

// HashToken creates a SHA256 hash of the token for secure storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
