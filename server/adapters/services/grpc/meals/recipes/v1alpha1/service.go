package v1alpha1

import (
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/fieldmasker"
	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
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
	RecipeNamer       namer.ReflectNamer `name:"v1alpha1RecipeNamer"`
	AccessNamer       namer.ReflectNamer `name:"v1alpha1RecipeAccessNamer"`
	UserNamer         namer.ReflectNamer `name:"v1alpha1UserNamer"`
	CircleNamer       namer.ReflectNamer `name:"v1alpha1CircleNamer"`
}

// NewRecipeService creates a new RecipeService.
func NewRecipeService(params NewRecipeServiceParams) (*RecipeService, error) {
	return &RecipeService{
		domain:            params.Domain,
		log:               params.Log,
		recipeFieldMasker: params.RecipeFieldMasker,
		recipeNamer:       params.RecipeNamer,
		userNamer:         params.UserNamer,
		accessNamer:       params.AccessNamer,
		circleNamer:       params.CircleNamer,
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
	userNamer         namer.ReflectNamer
	circleNamer       namer.ReflectNamer
	accessNamer       namer.ReflectNamer
}
