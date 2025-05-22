package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteRecipe -
func (s *RecipeService) DeleteRecipe(ctx context.Context, request *pb.DeleteRecipeRequest) (*pb.Recipe, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	parent, id, err := s.recipeNamer.Parse(request.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	if s.domain.AuthorizeRecipeParent(ctx, user, parent) != nil {
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
