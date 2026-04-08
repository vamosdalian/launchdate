package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vamosdalian/launchdate-backend/internal/util"
)

const (
	authorizationHeader = "Authorization"
	authorizationBearer = "Bearer"
	userIDKey           = "user_id"
	usernameKey         = "username"
	userRoleKey         = "user_role"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtManager *util.JWTManager, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token from Authorization header
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != authorizationBearer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// verify token
		claims, err := jwtManager.VerifyAccessToken(token)
		if err != nil {
			logger.WithError(err).Warn("failed to verify access token")

			var errorMsg string
			if err == util.ErrExpiredToken {
				errorMsg = "token expired"
			} else {
				errorMsg = "invalid token"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errorMsg,
			})
			c.Abort()
			return
		}

		// store user info in context
		c.Set(userIDKey, claims.UserID)
		c.Set(usernameKey, claims.Username)
		c.Set(userRoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole creates a middleware that checks if user has required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(userRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role not found",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, requiredRole := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
		c.Abort()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get(userIDKey); exists {
		return userID.(string)
	}
	return ""
}

// GetUsername retrieves username from context
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(usernameKey); exists {
		return username.(string)
	}
	return ""
}

// GetUserRole retrieves user role from context
func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get(userRoleKey); exists {
		return role.(string)
	}
	return ""
}
