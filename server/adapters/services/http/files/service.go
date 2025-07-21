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
	log               zerolog.Logger
	domain            domain.Domain
	recipeNamer       namer.ReflectNamer
	circleNamer       namer.ReflectNamer
	recipeAccessNamer namer.ReflectNamer

	userNamer namer.ReflectNamer
}

type NewServiceParams struct {
	fx.In

	Log               zerolog.Logger
	Domain            domain.Domain
	RecipeNamer       namer.ReflectNamer `name:"v1alpha1RecipeNamer"`
	CircleNamer       namer.ReflectNamer `name:"v1alpha1CircleNamer"`
	RecipeAccessNamer namer.ReflectNamer `name:"v1alpha1RecipeAccessNamer"`

	UserNamer namer.ReflectNamer `name:"v1alpha1UserNamer"`
}

func NewService(params NewServiceParams) (*Service, error) {
	return &Service{
		log:               params.Log,
		domain:            params.Domain,
		recipeNamer:       params.RecipeNamer,
		circleNamer:       params.CircleNamer,
		recipeAccessNamer: params.RecipeAccessNamer,
		userNamer:         params.UserNamer,
	}, nil
}

func (s *Service) Register(m *http.ServeMux) error {
	r := mux.NewRouter().StrictSlash(true)

	s.log.Info().Msg("Registering files service routes")
	r.HandleFunc("/meals/v1alpha1/{name:recipes/[0-9]+}/image", s.UploadRecipeImage).Methods(http.MethodPut)
	r.HandleFunc("/meals/v1alpha1/{name:circles/[0-9]*/recipes/[0-9]+}/image", s.UploadRecipeImage).Methods(http.MethodPut)
	r.HandleFunc("/meals/v1alpha1/{name:recipes/[0-9]+}/image:generate", s.GenerateRecipeImage).Methods(http.MethodGet)
	r.HandleFunc("/meals/v1alpha1/{name:circles/[0-9]*/recipes/[0-9]+}/image:generate", s.GenerateRecipeImage).Methods(http.MethodGet)

	r.HandleFunc("/circles/v1alpha1/{name:circles/[0-9]+}/image", s.UploadCircleImage).Methods(http.MethodPut)
	r.HandleFunc("/users/v1alpha1/{name:users/[0-9]+}/image", s.UploadUserImage).Methods(http.MethodPut)

	r.HandleFunc("/meals/v1alpha1/recipes:ocr", s.OCRRecipe).Methods(http.MethodPost)

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
