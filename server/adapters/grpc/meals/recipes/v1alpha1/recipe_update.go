package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/grpc/metadata"
	"github.com/jcfug8/daylear/server/core/errz"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// UpdateRecipe -
func (s *RecipeService) UpdateRecipe(ctx context.Context, request *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, errz.NewInvalidArgument("missing or invalid authorization token")
	}

	modelMask, err := s.recipeFieldMasker.GetWriteMask(request.GetUpdateMask())
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid field mask")
	}

	recipe := request.GetRecipe()

	if !s.recipeNamer.IsMatch(recipe.GetName()) {
		return nil, errz.NewInvalidArgument("invalid name: %s", recipe.GetName())
	}

	mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, recipe)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid request data")
	}

	if s.domain.AuthorizeParent(ctx, authToken, mRecipe.Parent) != nil {
		return nil, errz.NewPermissionDenied("user not authorized")
	}

	mRecipe, err = s.domain.UpdateRecipe(ctx, mRecipe, modelMask)
	if err != nil {
		return nil, errz.Sanitize(err)
	}

	recipe, err = convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, errz.NewInternal("unable to prepare response")
	}

	return recipe, nil
}
