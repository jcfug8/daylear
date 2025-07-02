package grpcgateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
)

type RequestIDKey string

// MiddlewareManager handles the registration and chaining of middleware
type MiddlewareManager struct {
	logger zerolog.Logger
	mux    *runtime.ServeMux
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager(logger zerolog.Logger, mux *runtime.ServeMux) *MiddlewareManager {
	return &MiddlewareManager{
		logger: logger,
		mux:    mux,
	}
}

// WithRequestID adds a request ID to the context
func (m *MiddlewareManager) WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx = context.WithValue(ctx, RequestIDKey("request_id"), requestID)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
