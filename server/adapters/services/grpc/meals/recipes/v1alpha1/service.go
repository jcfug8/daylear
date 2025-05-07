package v1alpha1

import (
	"sync/atomic"

	fieldValidator "github.com/jcfug8/daylear/server/adapters/services/grpc/fieldbehaviorvalidator"
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/fieldmasker"
	namer "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// NewRecipeServiceParams defines the dependencies for the RecipeService.
type NewRecipeServiceParams struct {
	fx.In

	Domain            domain.Domain
	FieldValidator    fieldValidator.FieldBehaviorValidator
	Log               zerolog.Logger
	RecipeFieldMasker fieldMasker.RecipeFieldMasker
	RecipeNamer       namer.RecipeNamer
}

// NewRecipeService creates a new RecipeService.
func NewRecipeService(params NewRecipeServiceParams) *RecipeService {

	return &RecipeService{
		domain:                 params.Domain,
		fieldBehaviorValidator: params.FieldValidator,
		log:                    params.Log,
		recipeFieldMasker:      params.RecipeFieldMasker,
		recipeNamer:            params.RecipeNamer,
	}
}

// RecipeService defines the grpc handlers for the RecipeService.
type RecipeService struct {
	pb.UnimplementedRecipeServiceServer
	domain                 domain.Domain
	fieldBehaviorValidator fieldValidator.FieldBehaviorValidator
	log                    zerolog.Logger
	recipeFieldMasker      fieldMasker.RecipeFieldMasker
	recipeNamer            namer.RecipeNamer
	registered             atomic.Bool
}
