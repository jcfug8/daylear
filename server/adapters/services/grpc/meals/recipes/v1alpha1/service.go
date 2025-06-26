package v1alpha1

import (
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/fieldmasker"
	namer "github.com/jcfug8/daylear/server/core/namer"
	circlePb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	userPb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// NewRecipeServiceParams defines the dependencies for the RecipeService.
type NewRecipeServiceParams struct {
	fx.In

	Domain            domain.Domain
	Log               zerolog.Logger
	RecipeFieldMasker fieldMasker.RecipeFieldMasker
}

// NewRecipeService creates a new RecipeService.
func NewRecipeService(params NewRecipeServiceParams) (*RecipeService, error) {
	recipeNamer, err := namer.NewReflectNamer[*pb.Recipe]()
	if err != nil {
		return nil, err
	}

	userNamer, err := namer.NewReflectNamer[*userPb.User]()
	if err != nil {
		return nil, err
	}

	circleNamer, err := namer.NewReflectNamer[*circlePb.Circle]()
	if err != nil {
		return nil, err
	}

	accessNamer, err := namer.NewReflectNamer[*pb.Access]()
	if err != nil {
		return nil, err
	}

	return &RecipeService{
		domain:            params.Domain,
		log:               params.Log,
		recipeFieldMasker: params.RecipeFieldMasker,
		recipeNamer:       recipeNamer,
		userNamer:         userNamer,
		circleNamer:       circleNamer,
		accessNamer:       accessNamer,
	}, nil
}

// RecipeService defines the grpc handlers for the RecipeService.
type RecipeService struct {
	pb.UnimplementedRecipeServiceServer
	pb.UnimplementedRecipeAccessServiceServer
	domain            domain.Domain
	log               zerolog.Logger
	recipeFieldMasker fieldMasker.RecipeFieldMasker
	recipeNamer       namer.ReflectNamer
	publicCircleNamer namer.ReflectNamer
	circleNamer       namer.ReflectNamer
	userNamer         namer.ReflectNamer
	publicUserNamer   namer.ReflectNamer

	accessNamer namer.ReflectNamer
}
