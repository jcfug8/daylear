package logutil

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/rs/zerolog"
)

// These should match the keys used in middleware and auth
const (
	requestIDKey = "req_id" // Used in context by middleware
)

// ExtractLogFieldsFromContext extracts request ID, user ID, and circle ID from context.
func ExtractLogFieldsFromContext(ctx context.Context) (requestID string, userID int64, circleID int64) {
	// Request ID
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			requestID = s
		}
	}

	// User/Circle (via AuthTokenMiddleware)
	auth, err := headers.ParseAuthData(ctx)
	if err == nil {
		userID = auth.UserId
		circleID = auth.CircleId
	}

	return
}

// EnrichLoggerWithContext returns a logger with request ID, user ID, and circle ID fields attached.
func EnrichLoggerWithContext(log zerolog.Logger, ctx context.Context) zerolog.Logger {
	requestID, userID, circleID := ExtractLogFieldsFromContext(ctx)
	return log.With().
		Str(requestIDKey, requestID).
		Int64("user_id", userID).
		Int64("circle_id", circleID).
		Logger()
}
