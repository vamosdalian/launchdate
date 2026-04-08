package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/auth"
)

const (
	refreshTokenCookieName = "refreshToken"
	cookieMaxAge           = 7 * 24 * 60 * 60 // 7 days in seconds
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService  *auth.AuthService
	logger       *logrus.Logger
	isProduction bool
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.AuthService, logger *logrus.Logger, isProduction bool) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		logger:       logger,
		isProduction: isProduction,
	}
}

// Bootstrap handles initial user creation if no user exists
func (h *AuthHandler) Bootstrap(c *gin.Context) {
	// 检查是否已有用户
	hasUser, err := h.authService.HasAnyUser(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("bootstrap: failed to check user count")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if hasUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "bootstrap disabled: user already exists"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.authService.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		h.logger.WithError(err).Error("bootstrap: failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bootstrap success", "user": user})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// get user agent and IP
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	// authenticate user
	response, refreshToken, err := h.authService.Login(c.Request.Context(), req.Username, req.Password, userAgent, ipAddress)
	if err != nil {
		h.logger.WithError(err).Warn("login failed")

		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid username or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// set refresh token in HttpOnly cookie
	h.setRefreshTokenCookie(c, refreshToken)

	c.JSON(http.StatusOK, response)
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// get refresh token from cookie
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing refresh token",
		})
		return
	}

	// refresh access token
	response, newRefreshToken, err := h.authService.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		h.logger.WithError(err).Warn("refresh token failed")

		if errors.Is(err, auth.ErrInvalidToken) {
			// clear invalid cookie
			h.clearRefreshTokenCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired refresh token",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// set new refresh token in cookie (token rotation)
	h.setRefreshTokenCookie(c, newRefreshToken)

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// get refresh token from cookie
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		// if no cookie, consider it already logged out
		h.clearRefreshTokenCookie(c)
		c.JSON(http.StatusOK, gin.H{
			"message": "logged out successfully",
		})
		return
	}

	// revoke refresh token
	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		h.logger.WithError(err).Warn("logout failed")
		// even if revocation fails, clear the cookie
	}

	// clear cookie
	h.clearRefreshTokenCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

// Me retrieves current user information using refresh token cookie
func (h *AuthHandler) Me(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authService.GetUserByRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrUserNotFound) {
			h.clearRefreshTokenCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		h.logger.WithError(err).Error("me: failed to load user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// setRefreshTokenCookie sets the refresh token in HttpOnly cookie
func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   cookieMaxAge,
		Expires:  time.Now().Add(time.Duration(cookieMaxAge) * time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

// clearRefreshTokenCookie clears the refresh token cookie
func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	sameSiteMode := http.SameSiteLaxMode
	if !h.isProduction {
		sameSiteMode = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: sameSiteMode,
	})
}
