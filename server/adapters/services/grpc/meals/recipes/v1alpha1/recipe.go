package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/aip/pagination"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	recipeMaxPageSize     int32 = 1000
	recipeDefaultPageSize int32 = 100
)

// CreateRecipe -
func (s *RecipeService) CreateRecipe(ctx context.Context, request *pb.CreateRecipeRequest) (response *pb.Recipe, err error) {
	authAccount, err := headers.ParseAuthData(ctx, s.recipeNamer)
	if err != nil {
		return nil, err
	}

	// check field behavior
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err = fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	// convert proto to model
	pbRecipe := request.GetRecipe()
	pbRecipe.Name = ""

	mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, pbRecipe)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// create recipe
	mRecipe, err = s.domain.CreateRecipe(ctx, authAccount, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbRecipe, err = convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	fieldbehavior.ClearFields(pbRecipe, annotations.FieldBehavior_INPUT_ONLY)

	return pbRecipe, nil
}

// DeleteRecipe -
func (s *RecipeService) DeleteRecipe(ctx context.Context, request *pb.DeleteRecipeRequest) (*pb.Recipe, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.recipeNamer)
	if err != nil {
		return nil, err
	}

	mRecipe := model.Recipe{}
	_, err = s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe, err = s.domain.DeleteRecipe(ctx, authAccount, mRecipe.Id)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to delete recipe")
		return nil, status.Errorf(codes.Internal, "unable to delete recipe")
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}

// GetRecipe -
func (s *RecipeService) GetRecipe(ctx context.Context, request *pb.GetRecipeRequest) (*pb.Recipe, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.recipeNamer)
	if err != nil {
		return nil, err
	}

	fieldMask := s.recipeFieldMasker.GetFieldMaskFromCtx(ctx)

	readMask, err := s.recipeFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	mRecipe := model.Recipe{}
	_, err = s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe, err = s.domain.GetRecipe(ctx, authAccount, mRecipe.Id, readMask)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to get recipe")
		return nil, status.Errorf(codes.Internal, "unable to get recipe")
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}

// Listrecipes -
func (s *RecipeService) ListRecipes(ctx context.Context, request *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.recipeNamer)
	if err != nil {
		return nil, err
	}

	fieldMask := s.recipeFieldMasker.GetFieldMaskFromCtx(ctx)

	readMask, err := s.recipeFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if request.GetPageSize() == 0 {
		request.PageSize = recipeDefaultPageSize
	}
	request.PageSize = min(request.PageSize, recipeMaxPageSize)

	res, err := s.domain.ListRecipes(ctx, authAccount, request.GetPageSize(), pageToken.Offset, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	recipes, err := convert.RecipeListToProto(s.recipeNamer, res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return &pb.ListRecipesResponse{
		NextPageToken: pageToken.Next(request).String(),
		Recipes:       recipes,
	}, nil
}
