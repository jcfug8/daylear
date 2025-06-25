package files

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/openapi"
)

type Service struct {
}

func NewService() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Register(m *http.ServeMux) error {
	r := mux.NewRouter().StrictSlash(true)

	r.PathPrefix("/").Handler(http.StripPrefix("/openapi/", http.FileServer(http.FS(openapi.FS))))
	m.Handle("/openapi/", r)
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "openapi-service"
}
