package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
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

	mRecipe := model.Recipe{}
	nameIndex, err := s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe.Parent.UserId = user.Id.UserId

	mRecipe, err = s.domain.DeleteRecipe(ctx, mRecipe.Parent, mRecipe.Id)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to delete recipe")
		return nil, status.Errorf(codes.Internal, "unable to delete recipe")
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe, nameIndex)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}
