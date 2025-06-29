package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	recipeMaxPageSize     int32 = 1000
	recipeDefaultPageSize int32 = 100
)

// CreateRecipe -
func (s *RecipeService) CreateRecipe(ctx context.Context, request *pb.CreateRecipeRequest) (response *pb.Recipe, err error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		return nil, err
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
	grpc.ProcessResponseFieldBehavior(pbRecipe)

	return pbRecipe, nil
}

// DeleteRecipe -
func (s *RecipeService) DeleteRecipe(ctx context.Context, request *pb.DeleteRecipeRequest) (*pb.Recipe, error) {
	authAccount, err := headers.ParseAuthData(ctx)
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}

// GetRecipe -
func (s *RecipeService) GetRecipe(ctx context.Context, request *pb.GetRecipeRequest) (*pb.Recipe, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	fieldMask := s.recipeFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.recipeFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mRecipe := model.Recipe{}
	_, err = s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe, err = s.domain.GetRecipe(ctx, authAccount, mRecipe.Id, readMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return pbRecipe, nil
}

// UpdateRecipe -
func (s *RecipeService) UpdateRecipe(ctx context.Context, request *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	recipeProto := request.GetRecipe()
	var mRecipe model.Recipe
	_, err = s.recipeNamer.Parse(recipeProto.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", recipeProto.GetName())
	}

	fieldMask := request.GetUpdateMask()
	updateMask, err := s.recipeFieldMasker.GetWriteMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mRecipe, err = convert.ProtoToRecipe(s.recipeNamer, recipeProto)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mRecipe, err = s.domain.UpdateRecipe(ctx, authAccount, mRecipe, updateMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	recipeProto, err = convert.RecipeToProto(s.recipeNamer, mRecipe)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return recipeProto, nil
}

// ListRecipes -
func (s *RecipeService) ListRecipes(ctx context.Context, request *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: recipeDefaultPageSize,
		MaxPageSize:     recipeMaxPageSize,
	})
	if err != nil {
		return nil, err
	}
	request.PageSize = pageSize

	res, err := s.domain.ListRecipes(ctx, authAccount, request.GetPageSize(), pageToken.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	recipes, err := convert.RecipeListToProto(s.recipeNamer, res)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	response := &pb.ListRecipesResponse{
		Recipes: recipes,
	}

	if len(recipes) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	return response, nil
}
