package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	recipeMaxPageSize     int32 = 1000
	recipeDefaultPageSize int32 = 100
)

var recipeFieldMap = map[string][]string{
	"name":              {model.RecipeField_Parent, model.RecipeField_Id},
	"title":             {model.RecipeField_Title},
	"description":       {model.RecipeField_Description},
	"directions":        {model.RecipeField_Directions},
	"ingredient_groups": {model.RecipeField_IngredientGroups},
	"image_uri":         {model.RecipeField_ImageURI},
	"visibility":        {model.RecipeField_VisibilityLevel},
	"citation":          {model.RecipeField_Citation},
	"cook_duration":     {model.RecipeField_CookDurationSeconds},
	"prep_duration":     {model.RecipeField_PrepDurationSeconds},
	"total_duration":    {model.RecipeField_TotalDurationSeconds},
	"cooking_method":    {model.RecipeField_CookingMethod},
	"categories":        {model.RecipeField_Categories},
	"yield_amount":      {model.RecipeField_YieldAmount},
	"cuisines":          {model.RecipeField_Cuisines},
	"create_time":       {model.RecipeField_CreateTime},
	"update_time":       {model.RecipeField_UpdateTime},

	"recipe_access.name":             {model.RecipeField_RecipeAccess},
	"recipe_access.permission_level": {model.RecipeAccessField_PermissionLevel},
	"recipe_access.state":            {model.RecipeAccessField_State},
}

// CreateRecipe -
func (s *RecipeService) CreateRecipe(ctx context.Context, request *pb.CreateRecipeRequest) (response *pb.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to process request field behavior")
		return nil, err
	}

	// convert proto to model
	pbRecipe := request.GetRecipe()
	pbRecipe.Name = ""
	_, mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, pbRecipe)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	nameIndex, err := s.recipeNamer.ParseParent(request.GetParent(), &mRecipe.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create recipe
	mRecipe, err = s.domain.CreateRecipe(ctx, authAccount, mRecipe)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbRecipe, err = convert.RecipeToProto(s.recipeNamer, s.accessNamer, mRecipe, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbRecipe)
	log.Info().Msg("gRPC CreateRecipe success")
	return pbRecipe, nil
}

// DeleteRecipe -
func (s *RecipeService) DeleteRecipe(ctx context.Context, request *pb.DeleteRecipeRequest) (*pb.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mRecipe := model.Recipe{}
	nameIndex, err := s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe, err = s.domain.DeleteRecipe(ctx, authAccount, mRecipe.Parent, mRecipe.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, s.accessNamer, mRecipe, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC DeleteRecipe success")
	return pbRecipe, nil
}

// GetRecipe -
func (s *RecipeService) GetRecipe(ctx context.Context, request *pb.GetRecipeRequest) (*pb.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mRecipe := model.Recipe{}
	nameIndex, err := s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe, err = s.domain.GetRecipe(ctx, authAccount, mRecipe.Parent, mRecipe.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbRecipe, err := convert.RecipeToProto(s.recipeNamer, s.accessNamer, mRecipe, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC GetRecipe success")
	return pbRecipe, nil
}

// UpdateRecipe -
func (s *RecipeService) UpdateRecipe(ctx context.Context, request *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	fieldMask := request.GetUpdateMask()
	updateMask := s.recipeFieldMasker.Convert(fieldMask.GetPaths())

	recipeProto := request.GetRecipe()
	nameIndex, mRecipe, err := convert.ProtoToRecipe(s.recipeNamer, recipeProto)
	if err != nil {
		log.Error().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.Internal, err.Error())
	}

	mRecipe, err = s.domain.UpdateRecipe(ctx, authAccount, mRecipe, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	recipeProto, err = convert.RecipeToProto(s.recipeNamer, s.accessNamer, mRecipe, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC UpdateRecipe success")
	return recipeProto, nil
}

// ListRecipes -
func (s *RecipeService) ListRecipes(ctx context.Context, request *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListRecipes called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mRecipeParent := model.RecipeParent{}
	_, err = s.recipeNamer.ParseParent(request.GetParent(), &mRecipeParent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: recipeDefaultPageSize,
		MaxPageSize:     recipeMaxPageSize,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to setup pagination")
		return nil, err
	}
	request.PageSize = pageSize

	res, err := s.domain.ListRecipes(ctx, authAccount, mRecipeParent, request.GetPageSize(), pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListRecipes failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	recipes := make([]*pb.Recipe, len(res))
	for i, recipe := range res {
		recipeProto, err := convert.RecipeToProto(s.recipeNamer, s.accessNamer, recipe)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		recipes[i] = recipeProto
	}

	// check field behavior
	for _, recipeProto := range recipes {
		grpc.ProcessResponseFieldBehavior(recipeProto)
	}

	response := &pb.ListRecipesResponse{
		Recipes: recipes,
	}

	if len(recipes) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListRecipes success")
	return response, nil
}

// ScrapeRecipe -
func (s *RecipeService) ScrapeRecipe(ctx context.Context, request *pb.ScrapeRecipeRequest) (*pb.ScrapeRecipeResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ScrapeRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to process request field behavior")
		return nil, err
	}

	// scrape the recipe
	recipe, err := s.domain.ScrapeRecipe(ctx, authAccount, request.GetUri())
	if err != nil {
		log.Error().Err(err).Msg("domain.ScrapeRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert the recipe to a proto
	recipeProto, err := convert.RecipeToProto(s.recipeNamer, s.accessNamer, recipe)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC ScrapeRecipe success")
	return &pb.ScrapeRecipeResponse{
		Recipe: recipeProto,
	}, nil
}
