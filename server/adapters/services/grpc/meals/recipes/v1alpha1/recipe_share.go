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

	mRecipe := model.Recipe{}
	_, err := s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe.Parent.UserId = user.Id.UserId

	parents := make([]model.RecipeParent, 0, len(request.GetRecipients()))
	for _, recipient := range request.GetRecipients() {
		recipientRecipeParent := model.RecipeParent{}
		_, err := s.publicUserNamer.ParseParent(recipient, &recipientRecipeParent)
		if err != nil {
			_, err := s.publicCircleNamer.ParseParent(recipient, &recipientRecipeParent)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", recipient)
			}
		}
		parents = append(parents, recipientRecipeParent)
	}

	err = s.domain.ShareRecipe(ctx, mRecipe.Parent, parents, mRecipe.Id, request.GetPermission())
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to share recipe")
		return nil, status.Errorf(codes.Internal, "unable to share recipe")
	}

	return &pb.ShareRecipeResponse{}, nil
}

// UnshareRecipe removes sharing permissions for a recipe.
func (s *RecipeService) UnshareRecipe(ctx context.Context, request *pb.UnshareRecipeRequest) (*pb.UnshareRecipeResponse, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	mRecipe := model.Recipe{}
	_, err := s.recipeNamer.Parse(request.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mRecipe.Parent.UserId = user.Id.UserId

	// Parse recipients
	parents := make([]model.RecipeParent, 0, len(request.GetRecipients()))
	for _, recipient := range request.GetRecipients() {
		recipientRecipeParent := model.RecipeParent{}
		_, err := s.userNamer.ParseParent(recipient, &recipientRecipeParent)
		if err != nil {
			_, err := s.publicCircleNamer.ParseParent(recipient, &recipientRecipeParent)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", recipient)
			}
		}
		parents = append(parents, recipientRecipeParent)
	}

	err = s.domain.UnshareRecipe(ctx, mRecipe.Parent, parents, mRecipe.Id)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to unshare recipe")
		return nil, status.Errorf(codes.Internal, "unable to unshare recipe")
	}

	return &pb.UnshareRecipeResponse{}, nil
}
