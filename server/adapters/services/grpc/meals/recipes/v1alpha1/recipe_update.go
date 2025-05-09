package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/metadata"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateRecipe -
func (s *RecipeService) UpdateRecipe(ctx context.Context, request *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing or invalid authorization token")
	}

	modelMask, err := s.recipeFieldMasker.GetWriteMask(request.GetUpdateMask())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	recipe := request.GetRecipe()

	if !s.recipeNamer.IsMatch(recipe.GetName()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %s", recipe.GetName())
	}

	mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, recipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data")
	}

	if s.domain.AuthorizeRecipeParent(ctx, authToken, mRecipe.Parent) != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	mRecipe, err = s.domain.UpdateRecipe(ctx, mRecipe, modelMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	recipe, err = convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return recipe, nil
}
