package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/metadata"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteRecipe -
func (s *RecipeService) DeleteRecipe(ctx context.Context, request *pb.DeleteRecipeRequest) (*pb.Recipe, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing or invalid authorization token")
	}

	parent, id, err := s.recipeNamer.Parse(request.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	if s.domain.AuthorizeParent(ctx, authToken, parent) != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	mRecipe, err := s.domain.DeleteRecipe(ctx, parent, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	recipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return recipe, nil
}
