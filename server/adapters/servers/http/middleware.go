package http

import (
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// NewMiddlewareMux -
func NewMiddlewareMux(log zerolog.Logger, origins []string, mux *http.ServeMux) *MiddlewareMux {
	return &MiddlewareMux{
		ServeMux: mux,
		origins:  origins,
		log:      log,
	}
}

// MiddlewareMux -
type MiddlewareMux struct {
	*http.ServeMux
	origins []string
	log     zerolog.Logger
}

// ServeHTTP -
func (m *MiddlewareMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := alice.New()

	c = c.Append(hlog.NewHandler(m.log))

	c = c.Append(cors.New(cors.Options{
		AllowedOrigins:   m.origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{headers.AuthorizationHeaderKey, "Content-Type", headers.ActingAsCircleHeaderKey},
		AllowCredentials: true,
	}).Handler)

	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		base := path.Base(r.URL.Path)
		if base == "readyz" || base == "healthz" {
			return
		}

		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Str("content_type", r.Header.Get("Content-Type")).
			Str("content_length", r.Header.Get("Content-Length")).
			Msg("HTTP Request")
	}))

	c = c.Append(PrefixHandler("/api"))

	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	h := c.Then(m.ServeMux)
	h.ServeHTTP(w, r)
}

func PrefixHandler(prefix string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, prefix) {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
			} else {
				http.NotFound(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
