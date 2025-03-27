package ping

import (
	"log"
	"net/http"

	"github.com/jcfug8/daylear/server/ports/domain"
)

type Service struct {
	domain domain.Domain
}

func NewService(domain domain.Domain) *Service {
	return &Service{
		domain: domain,
	}
}

func (s *Service) Register(m *http.ServeMux) error {
	m.HandleFunc("/ping", s.Ping)
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "http-ping-service"
}

func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	err := s.domain.Ping()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
