package caldav

import "net/http"

func (s *Service) Event(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Event called")
}
