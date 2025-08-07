package openapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/openapi"
)

type Service struct {
}

func NewService() (*Service, error) {
	return &Service{}, nil
}

// fileServerWrapper wraps the http.FileServer to add logging
func fileServerWrapper(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("FileServer hit: %s", r.URL.Path)
		http.FileServer(fs).ServeHTTP(w, r)
	})
}

func (s *Service) Register(m *http.ServeMux) error {
	r := mux.NewRouter().StrictSlash(true)

	// Register specific routes first - they take precedence over PathPrefix
	r.HandleFunc("/openapi/specs.json", s.GetSpecs).Methods("GET")
	// Serve swagger files but exclude specs.json by using a more specific pattern
	r.PathPrefix("/").Handler(http.StripPrefix("/openapi/", fileServerWrapper(http.FS(openapi.FS)))).Methods("GET")
	m.Handle("/openapi/", r)
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "openapi-service"
}
