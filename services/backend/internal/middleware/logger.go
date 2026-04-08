package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger returns a gin middleware for logging
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		if path == "/api/v1/health" {
			return
		}

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Log errors
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.WithFields(logrus.Fields{
				"status":     c.Writer.Status(),
				"method":     c.Request.Method,
				"path":       path,
				"query":      query,
				"ip":         c.ClientIP(),
				"latency":    latency,
				"user-agent": c.Request.UserAgent(),
			}).Info("request processed")
		}
	}
}
