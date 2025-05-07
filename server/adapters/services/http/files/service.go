package files

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/namer"
	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/rs/zerolog"
)

type Service struct {
	log         zerolog.Logger
	domain      domain.Domain
	recipeNamer namer.RecipeNamer
}

func NewService(log zerolog.Logger, domain domain.Domain, recipeNamer namer.RecipeNamer) *Service {
	return &Service{
		log:         log,
		domain:      domain,
		recipeNamer: recipeNamer,
	}
}

func (s *Service) Register(m *http.ServeMux) error {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/meals/v1alpha1/{name:users/[0-9]+/recipes/[0-9]+}/image", s.UploadRecipeImage).Methods(http.MethodPut)

	// m.Handle("/files/", headers.NewAuthTokenMiddleware(s.domain)(http.StripPrefix("/files", r)))
	m.Handle("/files/", http.StripPrefix("/files", r))
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "http-auth"
}
