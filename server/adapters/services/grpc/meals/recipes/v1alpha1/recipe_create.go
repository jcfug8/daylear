package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/metadata"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateRecipe -
func (s *RecipeService) CreateRecipe(ctx context.Context, request *pb.CreateRecipeRequest) (response *pb.Recipe, err error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "missing or invalid authorization token")
	}

	recipe := request.GetRecipe()

	err = s.fieldBehaviorValidator.Validate(recipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	recipe.Name = ""

	mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, recipe)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	mRecipe.Parent, err = s.recipeNamer.ParseParent(request.GetParent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	if s.domain.AuthorizeParent(ctx, authToken, mRecipe.Parent) != nil {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	mRecipe, err = s.domain.CreateRecipe(ctx, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	recipe, err = convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return recipe, nil
}
