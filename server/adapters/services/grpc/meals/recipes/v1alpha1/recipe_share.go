package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShareRecipe -
func (s *RecipeService) ShareRecipe(ctx context.Context, request *pb.ShareRecipeRequest) (*pb.ShareRecipeResponse, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	parent, id, err := s.recipeNamer.Parse(request.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	parents := make([]model.RecipeParent, 0, len(request.GetRecipients()))
	for _, recipient := range request.GetRecipients() {
		parent, err := s.recipeNamer.ParseParent(recipient)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", recipient)
		}
		parents = append(parents, parent)
	}

	if s.domain.AuthorizeRecipeParent(ctx, user, parent) != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	err = s.domain.ShareRecipe(ctx, parent, parents, id, request.GetPermission())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ShareRecipeResponse{}, nil
}
