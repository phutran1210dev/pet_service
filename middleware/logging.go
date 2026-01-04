package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware logs all requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := uuid.New().String()

		// Set request ID in context
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Log after request
		duration := time.Since(startTime)
		log.Printf("[%s] %s %s - Status: %d - Duration: %v - RequestID: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
			requestID,
		)
	}
}
