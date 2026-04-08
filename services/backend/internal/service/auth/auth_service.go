package auth

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthService handles authentication business logic
type AuthService struct {
	db         *mongo.Database
	jwtManager *util.JWTManager
	logger     *logrus.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(db *mongo.Database, jwtManager *util.JWTManager, logger *logrus.Logger) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// HasAnyUser checks if any user exists in the database
func (s *AuthService) HasAnyUser(ctx context.Context) (bool, error) {
	count, err := s.db.Collection("users").CountDocuments(ctx, bson.M{})
	if err != nil {
		s.logger.WithError(err).Error("failed to count users")
		return false, err
	}
	return count > 0, nil
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*models.LoginResponse, string, error) {
	// find user by username
	var user models.User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", ErrInvalidCredentials
		}
		s.logger.WithError(err).Error("failed to find user")
		return nil, "", err
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// generate access token
	accessToken, err := s.jwtManager.GenerateAccessToken(&user)
	if err != nil {
		s.logger.WithError(err).Error("failed to generate access token")
		return nil, "", err
	}

	// generate refresh token
	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		s.logger.WithError(err).Error("failed to generate refresh token")
		return nil, "", err
	}

	// store refresh token in database
	tokenHash := util.HashToken(refreshToken)
	refreshTokenDoc := models.RefreshToken{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiration()),
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	_, err = s.db.Collection("refresh_tokens").InsertOne(ctx, refreshTokenDoc)
	if err != nil {
		s.logger.WithError(err).Error("failed to store refresh token")
		return nil, "", err
	}

	// sanitize user response (remove password hash)
	userResp := user
	userResp.PasswordHash = ""

	return &models.LoginResponse{
		AccessToken: accessToken,
		User:        &userResp,
	}, refreshToken, nil
}

// RefreshAccessToken validates refresh token and generates new access token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.RefreshResponse, string, error) {
	// verify refresh token format
	if err := s.jwtManager.VerifyRefreshToken(refreshToken); err != nil {
		return nil, "", ErrInvalidToken
	}

	// check if token exists and is valid in database
	tokenHash := util.HashToken(refreshToken)
	var storedToken models.RefreshToken
	err := s.db.Collection("refresh_tokens").FindOne(ctx, bson.M{
		"token_hash": tokenHash,
		"revoked_at": nil,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&storedToken)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", ErrInvalidToken
		}
		s.logger.WithError(err).Error("failed to find refresh token")
		return nil, "", err
	}

	// get user
	var user models.User
	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": storedToken.UserID}).Decode(&user)
	if err != nil {
		s.logger.WithError(err).Error("failed to find user")
		return nil, "", ErrUserNotFound
	}

	// update last used timestamp
	now := time.Now()
	_, err = s.db.Collection("refresh_tokens").UpdateOne(
		ctx,
		bson.M{"_id": storedToken.ID},
		bson.M{"$set": bson.M{"last_used_at": now}},
	)
	if err != nil {
		s.logger.WithError(err).Warn("failed to update last_used_at")
	}

	// generate new access token
	accessToken, err := s.jwtManager.GenerateAccessToken(&user)
	if err != nil {
		s.logger.WithError(err).Error("failed to generate access token")
		return nil, "", err
	}

	// optional: implement refresh token rotation
	// generate new refresh token
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		s.logger.WithError(err).Error("failed to generate new refresh token")
		return nil, "", err
	}

	// revoke old token
	_, err = s.db.Collection("refresh_tokens").UpdateOne(
		ctx,
		bson.M{"_id": storedToken.ID},
		bson.M{"$set": bson.M{"revoked_at": now}},
	)
	if err != nil {
		s.logger.WithError(err).Warn("failed to revoke old refresh token")
	}

	// store new refresh token
	newTokenHash := util.HashToken(newRefreshToken)
	newTokenDoc := models.RefreshToken{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiration()),
		CreatedAt: time.Now(),
		UserAgent: storedToken.UserAgent,
		IPAddress: storedToken.IPAddress,
	}

	_, err = s.db.Collection("refresh_tokens").InsertOne(ctx, newTokenDoc)
	if err != nil {
		s.logger.WithError(err).Error("failed to store new refresh token")
		return nil, "", err
	}

	return &models.RefreshResponse{
		AccessToken: accessToken,
	}, newRefreshToken, nil
}

// Logout revokes the refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := util.HashToken(refreshToken)

	result, err := s.db.Collection("refresh_tokens").UpdateOne(
		ctx,
		bson.M{
			"token_hash": tokenHash,
			"revoked_at": nil,
		},
		bson.M{"$set": bson.M{"revoked_at": time.Now()}},
	)

	if err != nil {
		s.logger.WithError(err).Error("failed to revoke refresh token")
		return err
	}

	if result.MatchedCount == 0 {
		return ErrInvalidToken
	}

	return nil
}

// GetUserByRefreshToken validates refresh token and returns associated user
func (s *AuthService) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error) {
	if err := s.jwtManager.VerifyRefreshToken(refreshToken); err != nil {
		return nil, ErrInvalidToken
	}

	tokenHash := util.HashToken(refreshToken)
	var storedToken models.RefreshToken
	err := s.db.Collection("refresh_tokens").FindOne(ctx, bson.M{
		"token_hash": tokenHash,
		"revoked_at": nil,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&storedToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalidToken
		}
		s.logger.WithError(err).Error("failed to find refresh token")
		return nil, err
	}

	var user models.User
	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": storedToken.UserID}).Decode(&user)
	if err != nil {
		s.logger.WithError(err).Error("failed to find user")
		return nil, ErrUserNotFound
	}

	now := time.Now()
	_, err = s.db.Collection("refresh_tokens").UpdateOne(
		ctx,
		bson.M{"_id": storedToken.ID},
		bson.M{"$set": bson.M{"last_used_at": now}},
	)
	if err != nil {
		s.logger.WithError(err).Warn("failed to update last_used_at")
	}

	user.PasswordHash = ""
	return &user, nil
}

// CreateUser creates a new user (for admin or registration)
func (s *AuthService) CreateUser(ctx context.Context, username, email, password, role string) (*models.User, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("failed to hash password")
		return nil, err
	}

	user := models.User{
		ID:           primitive.NewObjectID(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("failed to create user")
		return nil, err
	}

	user.PasswordHash = "" // sanitize
	return &user, nil
}

// CleanupExpiredTokens removes expired refresh tokens from database
func (s *AuthService) CleanupExpiredTokens(ctx context.Context) error {
	_, err := s.db.Collection("refresh_tokens").DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to cleanup expired tokens")
		return err
	}
	return nil
}
