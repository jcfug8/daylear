package caldav

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type Service struct {
	log zerolog.Logger

	domain domain.Domain
}

type NewServiceParams struct {
	fx.In

	Log    zerolog.Logger
	Domain domain.Domain
}

func NewService(params NewServiceParams) (*Service, error) {
	s := &Service{
		log:    params.Log,
		domain: params.Domain,
	}

	return s, nil
}

// TODO: make sure path username/user id param matches the auth account user id on headers
func (s *Service) Register(m *http.ServeMux) error {
	s.log.Info().Msg("Registering caldav service routes")

	wellKnownGMux := mux.NewRouter()
	wellKnownGMux.HandleFunc("/.well-known/caldav", s.WellKnown).Methods("GET", "OPTIONS", "PROPFIND")
	m.Handle("/.well-known/caldav", headers.NewBasicAuthMiddleware(s.domain)(wellKnownGMux))

	gmux := mux.NewRouter()

	// These return the principal collection set but they aren't needed right now
	// gmux.HandleFunc("/caldav", s.Root).Methods("OPTIONS", "PROPFIND")
	// gmux.HandleFunc("/caldav/", s.Root).Methods("OPTIONS", "PROPFIND")

	// These currently just redirect to the current user principal
	gmux.HandleFunc("/caldav/.well-known/caldav", s.WellKnown).Methods("OPTIONS", "PROPFIND")
	gmux.HandleFunc("/caldav/.well-known/caldav/", s.WellKnown).Methods("OPTIONS", "PROPFIND")

	// These return the principal collection set but they aren't needed right now
	// gmux.HandleFunc("/caldav/principals", s.Principals).Methods("PROPFIND", "OPTIONS")
	// gmux.HandleFunc("/caldav/principals/", s.Principals).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}", s.UserPrincipal).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}/", s.UserPrincipal).Methods("PROPFIND", "OPTIONS")

	// These return the calendar home set but they aren't needed right now
	// gmux.HandleFunc("/caldav/principals/{userID}/calendar-home-set", s.CalendarHomeSet).Methods("PROPFIND", "OPTIONS")
	// gmux.HandleFunc("/caldav/principals/{userID}/calendar-home-set/", s.CalendarHomeSet).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}/calendars", s.Calendars).Methods("PROPFIND", "OPTIONS")
	gmux.HandleFunc("/caldav/principals/{userID}/calendars/", s.Calendars).Methods("PROPFIND", "OPTIONS")

	// Add calendar objects endpoints for individual calendars and events
	gmux.HandleFunc("/caldav/principals/{userID}/calendars/{calendarID}", s.Calendar).Methods("OPTIONS", "PROPFIND", "REPORT")
	gmux.HandleFunc("/caldav/principals/{userID}/calendars/{calendarID}/", s.Calendar).Methods("OPTIONS", "PROPFIND", "REPORT")
	// gmux.HandleFunc("/caldav/principals/{userID}/calendars/{calendarID}/{eventID}.ics", s.Event).Methods("OPTIONS", "GET", "PUT", "DELETE")

	m.Handle("/caldav/", headers.NewBasicAuthMiddleware(s.domain)(gmux))

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "caldav-service"
}
