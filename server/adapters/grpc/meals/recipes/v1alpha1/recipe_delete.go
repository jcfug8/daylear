package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/grpc/metadata"
	"github.com/jcfug8/daylear/server/core/errz"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// DeleteRecipe -
func (s *RecipeService) DeleteRecipe(ctx context.Context, request *pb.DeleteRecipeRequest) (*pb.Recipe, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, errz.NewInvalidArgument("missing or invalid authorization token")
	}

	parent, id, err := s.recipeNamer.Parse(request.GetName())
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid name: %v", request.GetName())
	}

	if s.domain.AuthorizeParent(ctx, authToken, parent) != nil {
		return nil, errz.NewPermissionDenied("user not authorized")
	}

	mRecipe, err := s.domain.DeleteRecipe(ctx, parent, id)
	if err != nil {
		return nil, errz.Sanitize(err)
	}

	recipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, errz.NewInternal("unable to prepare response")
	}

	return recipe, nil
}
