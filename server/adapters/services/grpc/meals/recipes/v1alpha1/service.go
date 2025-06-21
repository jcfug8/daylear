package v1alpha1

import (
	fieldValidator "github.com/jcfug8/daylear/server/adapters/services/grpc/fieldbehaviorvalidator"
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
	FieldValidator    fieldValidator.FieldBehaviorValidator
	Log               zerolog.Logger
	RecipeFieldMasker fieldMasker.RecipeFieldMasker
}

// NewRecipeService creates a new RecipeService.
func NewRecipeService(params NewRecipeServiceParams) (*RecipeService, error) {
	recipeNamer, err := namer.NewReflectNamer[*pb.Recipe]()
	if err != nil {
		return nil, err
	}

	recipeNamer_PublicCircle, err := namer.NewReflectNamer[*circlePb.PublicCircle]()
	if err != nil {
		return nil, err
	}

	recipeNamer_User, err := namer.NewReflectNamer[*userPb.User]()
	if err != nil {
		return nil, err
	}

	recipeNamer_Circle, err := namer.NewReflectNamer[*circlePb.Circle]()
	if err != nil {
		return nil, err
	}

	recipeNamer_PublicUser, err := namer.NewReflectNamer[*userPb.PublicUser]()
	if err != nil {
		return nil, err
	}

	accessNamer, err := namer.NewReflectNamer[*pb.Access]()
	if err != nil {
		return nil, err
	}

	return &RecipeService{
		domain:                 params.Domain,
		fieldBehaviorValidator: params.FieldValidator,
		log:                    params.Log,
		recipeFieldMasker:      params.RecipeFieldMasker,
		recipeNamer:            recipeNamer,
		publicCircleNamer:      recipeNamer_PublicCircle,
		userNamer:              recipeNamer_User,
		circleNamer:            recipeNamer_Circle,
		publicUserNamer:        recipeNamer_PublicUser,
		accessNamer:            accessNamer,
	}, nil
}

// RecipeService defines the grpc handlers for the RecipeService.
type RecipeService struct {
	pb.UnimplementedRecipeServiceServer
	pb.UnimplementedRecipeRecipientsServiceServer
	pb.UnimplementedRecipeAccessServiceServer
	domain                 domain.Domain
	fieldBehaviorValidator fieldValidator.FieldBehaviorValidator
	log                    zerolog.Logger
	recipeFieldMasker      fieldMasker.RecipeFieldMasker
	recipeNamer            namer.ReflectNamer
	publicCircleNamer      namer.ReflectNamer
	circleNamer            namer.ReflectNamer
	userNamer              namer.ReflectNamer
	publicUserNamer        namer.ReflectNamer

	accessNamer namer.ReflectNamer
}
