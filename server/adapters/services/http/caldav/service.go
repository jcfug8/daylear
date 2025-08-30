package caldav

import (
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type Service struct {
	log zerolog.Logger

	domain  domain.Domain
	apiPath string
}

type NewServiceParams struct {
	fx.In

	Log    zerolog.Logger
	Domain domain.Domain
	Config config.Client
}

func NewService(params NewServiceParams) (*Service, error) {
	apiDomainConfig := params.Config.GetConfig()["apidomain"].(map[string]interface{})

	apiPath, _ := apiDomainConfig["path"].(string)
	if apiPath != "" && !strings.HasPrefix(apiPath, "/") {
		apiPath = "/" + apiPath
	}

	s := &Service{
		log:     params.Log,
		domain:  params.Domain,
		apiPath: apiPath,
	}

	return s, nil
}

// TODO: make sure path username/user id param matches the auth account user id on headers
func (s *Service) Register(m *http.ServeMux) error {
	s.log.Info().Msg("Registering caldav service routes")

	wellKnownGMux := mux.NewRouter()
	wellKnownGMux.HandleFunc("/.well-known/caldav", s.WellKnown).Methods("GET", "OPTIONS", "PROPFIND")
	m.Handle("/.well-known/caldav", headers.NewBasicAuthMiddleware(s.domain)(wellKnownGMux))
	wellKnownGMux.HandleFunc("/.well-known/caldav/", s.WellKnown).Methods("GET", "OPTIONS", "PROPFIND")
	m.Handle("/.well-known/caldav/", headers.NewBasicAuthMiddleware(s.domain)(wellKnownGMux))

	gmux := mux.NewRouter()

	gmux.HandleFunc("/caldav", s.Root).Methods("OPTIONS", "PROPFIND")
	gmux.HandleFunc("/caldav/", s.Root).Methods("OPTIONS", "PROPFIND")

	// These currently just redirect to the current user principal
	gmux.HandleFunc("/caldav/.well-known/caldav", s.WellKnown).Methods("OPTIONS", "PROPFIND")
	gmux.HandleFunc("/caldav/.well-known/caldav/", s.WellKnown).Methods("OPTIONS", "PROPFIND")

	gmux.HandleFunc("/caldav/principals", s.Principals).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/", s.Principals).Methods("PROPFIND", "OPTIONS")

	gmux.HandleFunc("/caldav/principals/{userID}", s.UserPrincipal).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}/", s.UserPrincipal).Methods("PROPFIND", "OPTIONS")

	gmux.HandleFunc("/caldav/principals/{userID}/calendar-home-set", s.CalendarHomeSet).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}/calendar-home-set/", s.CalendarHomeSet).Methods("PROPFIND", "OPTIONS")

	gmux.HandleFunc("/caldav/principals/{userID}/calendars", s.Calendars).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}/calendars/", s.Calendars).Methods("PROPFIND", "OPTIONS")

	// Add calendar objects endpoints for individual calendars and events
	gmux.HandleFunc("/caldav/principals/{userID}/calendars/{calendarID}", s.Calendar).Methods("OPTIONS", "PROPFIND", "REPORT", "GET")
	gmux.HandleFunc("/caldav/principals/{userID}/calendars/{calendarID}/", s.Calendar).Methods("OPTIONS", "PROPFIND", "REPORT", "GET")

	gmux.HandleFunc("/caldav/principals/{userID}/calendars/{calendarID}/events/{eventID}.ics", s.Event).Methods("OPTIONS", "GET")

	m.Handle("/caldav", headers.NewBasicAuthMiddleware(s.domain)(gmux))
	m.Handle("/caldav/", headers.NewBasicAuthMiddleware(s.domain)(gmux))

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "caldav-service"
}

func (s *Service) NewResponseHref(p string) ResponseHref {
	href := path.Join(s.apiPath, p)
	if strings.HasSuffix(p, "/") && !strings.HasSuffix(href, "/") {
		href = href + "/"
	}
	return ResponseHref{
		Href: href,
	}
}

func (s *Service) NewResponseHrefPointer(p string) *ResponseHref {
	responseHref := s.NewResponseHref(p)
	return &responseHref
}
