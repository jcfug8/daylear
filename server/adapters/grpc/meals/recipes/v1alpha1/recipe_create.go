package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/grpc/metadata"
	"github.com/jcfug8/daylear/server/core/errz"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// CreateRecipe -
func (s *RecipeService) CreateRecipe(ctx context.Context, request *pb.CreateRecipeRequest) (response *pb.Recipe, err error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, errz.NewInvalidArgument("missing or invalid authorization token")
	}

	recipe := request.GetRecipe()

	err = s.fieldBehaviorValidator.Validate(recipe)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid request data: %v", err)
	}

	recipe.Name = ""

	mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, recipe)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, errz.NewInvalidArgument("invalid request data")
	}

	mRecipe.Parent, err = s.recipeNamer.ParseParent(request.GetParent())
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid parent: %v", request.GetParent())
	}

	if s.domain.AuthorizeParent(ctx, authToken, mRecipe.Parent) != nil {
		return nil, errz.NewPermissionDenied("user not authorized")
	}

	mRecipe, err = s.domain.CreateRecipe(ctx, mRecipe)
	if err != nil {
		return nil, errz.Sanitize(err)
	}

	recipe, err = convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, errz.NewInternal("unable to prepare response")
	}

	return recipe, nil
}
