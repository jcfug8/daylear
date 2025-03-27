package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/grpc/metadata"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// ShareRecipe -
func (s *RecipeService) ShareRecipe(ctx context.Context, request *pb.ShareRecipeRequest) (*pb.ShareRecipeResponse, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, errz.NewInvalidArgument("missing or invalid authorization token")
	}

	parent, id, err := s.recipeNamer.Parse(request.GetName())
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid name: %v", request.GetName())
	}

	parents := make([]model.RecipeParent, 0, len(request.GetRecipients()))
	for _, recipient := range request.GetRecipients() {
		parent, err := s.recipeNamer.ParseParent(recipient)
		if err != nil {
			return nil, errz.NewInvalidArgument("invalid recipient: %v", recipient)
		}
		parents = append(parents, parent)
	}

	if s.domain.AuthorizeParent(ctx, authToken, parent) != nil {
		return nil, errz.NewPermissionDenied("user not authorized")
	}

	err = s.domain.ShareRecipe(ctx, parent, parents, id, request.GetPermission())
	if err != nil {
		return nil, errz.Sanitize(err)
	}

	return &pb.ShareRecipeResponse{}, nil
}
