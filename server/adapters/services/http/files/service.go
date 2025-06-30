package files

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/namer"
	"github.com/jcfug8/daylear/server/ports/domain"
	"go.uber.org/fx"

	"github.com/rs/zerolog"
)

type Service struct {
	log         zerolog.Logger
	domain      domain.Domain
	recipeNamer namer.ReflectNamer
}

type NewServiceParams struct {
	fx.In

	Log         zerolog.Logger
	Domain      domain.Domain
	RecipeNamer namer.ReflectNamer `name:"v1alpha1RecipeNamer"`
}

func NewService(params NewServiceParams) (*Service, error) {
	return &Service{
		log:         params.Log,
		domain:      params.Domain,
		recipeNamer: params.RecipeNamer,
	}, nil
}

func (s *Service) Register(m *http.ServeMux) error {
	r := mux.NewRouter().StrictSlash(true)

	s.log.Info().Msg("Registering files service routes")
	r.HandleFunc("/meals/v1alpha1/{name:recipes/[0-9]+}/image", s.UploadRecipeImage).Methods(http.MethodPut)

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
