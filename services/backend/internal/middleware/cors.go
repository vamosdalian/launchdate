package middleware

import (
	"github.com/gin-gonic/gin"
)

var allowedOrigins = map[string]bool{
	"https://launch-date.com":       true,
	"https://www.launch-date.com":   true,
	"https://admin.launch-date.com": true,
}

// CORS returns a gin middleware for handling CORS
func CORS(isprod bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if isprod {
			if allowedOrigins[origin] {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.AbortWithStatus(403)
				return
			}
		} else {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
