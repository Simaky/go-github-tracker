package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Simaky/go-github-tracker/backend/server/handlers"
)

const requestIDHeader = "X-Request-ID"

// RequestID assigns (or echoes) a request id, stores it in the gin context for
// handlers, and returns it in the response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(handlers.RequestIDKey, id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}

// RequestLogger logs only failed requests; we do not log info-level lines per
// successful request.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if status := c.Writer.Status(); status >= http.StatusInternalServerError {
			slog.Error("request failed",
				"request_id", handlers.RequestIDFrom(c),
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
			)
		}
	}
}
