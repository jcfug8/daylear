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
	userV1alpha1Service         userV1alpha1.UserServiceServer
	userSettingsV1alpha1Service userV1alpha1.UserSettingsServiceServer
	userAccessService           userV1alpha1.UserAccessServiceServer
	recipesV1alpha1Service      recipesV1alpha1.RecipeServiceServer
	recipeAccessService         recipesV1alpha1.RecipeAccessServiceServer
	circleV1alpha1Service       circleV1alpha1.CircleServiceServer
	circleAccessService         circleV1alpha1.CircleAccessServiceServer
	domain                      domain.Domain
}

type NewServiceParams struct {
	fx.In

	UserV1alpha1Service         userV1alpha1.UserServiceServer
	UserSettingsV1alpha1Service userV1alpha1.UserSettingsServiceServer
	UserAccessService           userV1alpha1.UserAccessServiceServer
	RecipesV1alpha1Service      recipesV1alpha1.RecipeServiceServer
	RecipeAccessService         recipesV1alpha1.RecipeAccessServiceServer
	CircleV1alpha1Service       circleV1alpha1.CircleServiceServer
	CircleAccessService         circleV1alpha1.CircleAccessServiceServer
	Domain                      domain.Domain
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userV1alpha1Service:         params.UserV1alpha1Service,
		userSettingsV1alpha1Service: params.UserSettingsV1alpha1Service,
		userAccessService:           params.UserAccessService,
		recipesV1alpha1Service:      params.RecipesV1alpha1Service,
		recipeAccessService:         params.RecipeAccessService,
		circleV1alpha1Service:       params.CircleV1alpha1Service,
		circleAccessService:         params.CircleAccessService,
		domain:                      params.Domain,
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
	err = userV1alpha1.RegisterUserSettingsServiceHandlerServer(ctx, mux, s.userSettingsV1alpha1Service)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
		return err
	}
	err = userV1alpha1.RegisterUserAccessServiceHandlerServer(ctx, mux, s.userAccessService)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
		return err
	}
	err = recipesV1alpha1.RegisterRecipeServiceHandlerServer(ctx, mux, s.recipesV1alpha1Service)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
	}
	err = recipesV1alpha1.RegisterRecipeAccessServiceHandlerServer(ctx, mux, s.recipeAccessService)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
	}
	err = circleV1alpha1.RegisterCircleServiceHandlerServer(ctx, mux, s.circleV1alpha1Service)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
	}
	err = circleV1alpha1.RegisterCircleAccessServiceHandlerServer(ctx, mux, s.circleAccessService)
	if err != nil {
		log.Printf("Failed to register gRPC gateway: %v", err)
	}

	m.Handle("/", headers.NewAuthTokenMiddleware(s.domain)(mux))
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "http-user"
}
