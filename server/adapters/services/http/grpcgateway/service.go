package grpcgateway

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	circleV1alpha1 "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	recipesV1alpha1 "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	userV1alpha1 "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/domain"
)

type Service struct {
	userV1alpha1Service    userV1alpha1.UserServiceServer
	recipesV1alpha1Service recipesV1alpha1.RecipeServiceServer
	circleV1alpha1Service  circleV1alpha1.CircleServiceServer
	domain                 domain.Domain
}

type NewServiceParams struct {
	fx.In

	UserV1alpha1Service    userV1alpha1.UserServiceServer
	RecipesV1alpha1Service recipesV1alpha1.RecipeServiceServer
	CircleV1alpha1Service  circleV1alpha1.CircleServiceServer
	Domain                 domain.Domain
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userV1alpha1Service:    params.UserV1alpha1Service,
		recipesV1alpha1Service: params.RecipesV1alpha1Service,
		circleV1alpha1Service:  params.CircleV1alpha1Service,
		domain:                 params.Domain,
	}
}

func (s *Service) Register(m *http.ServeMux) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := userV1alpha1.RegisterUserServiceHandlerServer(ctx, mux, s.userV1alpha1Service)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
		return err
	}
	err = recipesV1alpha1.RegisterRecipeServiceHandlerServer(ctx, mux, s.recipesV1alpha1Service)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
	}
	err = circleV1alpha1.RegisterCircleServiceHandlerServer(ctx, mux, s.circleV1alpha1Service)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
	}

	m.Handle("/", headers.NewAuthTokenMiddleware(s.domain)(mux))
	// m.Handle("/", mux)
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "http-user"
}
