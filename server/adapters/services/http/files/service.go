package files

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
)

type Service struct {
	log         zerolog.Logger
	domain      domain.Domain
	recipeNamer namer.ReflectNamer
}

func NewService(log zerolog.Logger, domain domain.Domain) (*Service, error) {
	recipeNamer, err := namer.NewReflectNamer[*pb.Recipe]()
	if err != nil {
		return nil, err
	}

	return &Service{
		log:         log,
		domain:      domain,
		recipeNamer: recipeNamer,
	}, nil
}

func (s *Service) Register(m *http.ServeMux) error {
	r := mux.NewRouter().StrictSlash(true)

	s.log.Info().Msg("Registering files service routes")
	r.HandleFunc("/meals/v1alpha1/{name:users/[0-9]+/recipes/[0-9]+}/image", s.UploadRecipeImage).Methods(http.MethodPut)

	s.log.Info().Msg("Mounting files service at /files/")
	m.Handle("/files/", headers.NewAuthTokenMiddleware(s.domain)(http.StripPrefix("/files", r)))
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "files-service"
}
