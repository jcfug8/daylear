package http

import (
	"net/http"
	"path"
	"time"

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
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler)

	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		base := path.Base(r.URL.Path)
		if base == "readyz" || base == "healthz" {
			return
		}

		hlog.FromRequest(r).Debug().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	h := c.Then(m.ServeMux)
	h.ServeHTTP(w, r)
}
