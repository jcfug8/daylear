package v1alpha1

import (
	"context"
	"fmt"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/pagination"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	recipeMaxPageSize     int32 = 1000
	recipeDefaultPageSize int32 = 100
)

// Listrecipes -
func (s *RecipeService) ListRecipes(ctx context.Context, request *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}
	fmt.Println(user)

	fieldMask := s.recipeFieldMasker.GetFieldMaskFromCtx(ctx)

	mRecipe := model.Recipe{}
	_, err := s.recipeNamer.ParseParent(request.GetParent(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unable to parse parent: %v", request.GetParent())
	}

	readMask, err := s.recipeFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, err := pagination.ParsePageToken[cmodel.Recipe](request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if pageToken.PageSize == 0 {
		pageToken.PageSize = recipeDefaultPageSize
	}
	pageToken.PageSize = min(pageToken.PageSize, recipeMaxPageSize)

	if s.domain.AuthorizeRecipeParent(ctx, user, mRecipe.Parent) != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	res, err := s.domain.ListRecipes(ctx, pageToken, mRecipe.Parent, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	recipes, err := convert.RecipeListToProto(s.recipeNamer, res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return &pb.ListRecipesResponse{
		NextPageToken: pagination.EncodePageToken(pageToken.Next(res)),
		Recipes:       recipes,
	}, nil
}
