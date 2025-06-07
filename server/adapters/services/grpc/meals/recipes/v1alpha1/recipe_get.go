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

// GetRecipe -
func (s *RecipeService) GetRecipe(ctx context.Context, request *pb.GetRecipeRequest) (*pb.Recipe, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	fieldMask := s.recipeFieldMasker.GetFieldMaskFromCtx(ctx)

	readMask, err := s.recipeFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	mRecipe := model.Recipe{}
	nameIndex, err := s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe.Parent.UserId = user.Id.UserId

	mRecipe, err = s.domain.GetRecipe(ctx, mRecipe.Parent, mRecipe.Id, readMask)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to get recipe")
		return nil, status.Errorf(codes.Internal, "unable to get recipe")
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe, nameIndex)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}
