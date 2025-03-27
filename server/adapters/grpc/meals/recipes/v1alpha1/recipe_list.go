package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/grpc/metadata"
	"github.com/jcfug8/daylear/server/adapters/grpc/pagination"
	"github.com/jcfug8/daylear/server/core/errz"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// IRIOMO:CUSTOM_CODE_SLOT_START recipeServiceListConstants

const (
	recipeMaxPageSize     int32 = 1000
	recipeDefaultPageSize int32 = 100
)

// IRIOMO:CUSTOM_CODE_SLOT_END

// Listrecipes -
func (s *RecipeService) ListRecipes(ctx context.Context, request *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, errz.NewInvalidArgument("missing or invalid authorization token")
	}

	fieldMask := s.recipeFieldMasker.GetFieldMaskFromCtx(ctx)

	parent, err := s.recipeNamer.ParseParent(request.GetParent())
	if err != nil {
		return nil, errz.NewInvalidArgument("unable to parse parent: %v", request.GetParent())
	}

	readMask, err := s.recipeFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid field mask")
	}

	pageToken, err := pagination.ParsePageToken[cmodel.Recipe](request)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid page token")
	}

	if pageToken.PageSize == 0 {
		pageToken.PageSize = recipeDefaultPageSize
	}
	pageToken.PageSize = min(pageToken.PageSize, recipeMaxPageSize)

	if s.domain.AuthorizeParent(ctx, authToken, parent) != nil {
		return nil, errz.NewPermissionDenied("user not authorized")
	}

	res, err := s.domain.ListRecipes(ctx, pageToken, parent, request.GetFilter(), readMask)
	if err != nil {
		return nil, errz.Sanitize(err)
	}

	recipes, err := convert.RecipeListToProto(s.recipeNamer, res)
	if err != nil {
		return nil, errz.NewInternal("unable to prepare response")
	}

	return &pb.ListRecipesResponse{
		NextPageToken: pagination.EncodePageToken(pageToken.Next(res)),
		Recipes:       recipes,
	}, nil
}
