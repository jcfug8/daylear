package grpcgateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

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

// WithCORS adds CORS middleware
func (m *MiddlewareManager) WithCORS() http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure this based on your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler(m.mux)
}

// WithLogging adds request logging middleware
func (m *MiddlewareManager) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			status:        http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the request details
		m.logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rw.status).
			Str("duration", time.Since(start).String()).
			Msg("HTTP Request")
	})
}

// WithRequestID adds a request ID to the context
func (m *MiddlewareManager) WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		ctx = context.WithValue(ctx, "request_id", requestID)
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