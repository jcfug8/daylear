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

// UpdateRecipe -
func (s *RecipeService) UpdateRecipe(ctx context.Context, request *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	modelMask, err := s.recipeFieldMasker.GetWriteMask(request.GetUpdateMask())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	pbRecipe := request.GetRecipe()

	mRecipe, nameIndex, err := convert.ProtoToRecipe(s.recipeNamer, pbRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data")
	}

	mRecipe.Parent.UserId = user.Id.UserId

	mRecipe, err = s.domain.UpdateRecipe(ctx, mRecipe, modelMask)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to update recipe")
		return nil, status.Errorf(codes.Internal, "unable to update recipe")
	}

	pbRecipe, err = convert.RecipeToProto(s.recipeNamer, mRecipe, nameIndex)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}
