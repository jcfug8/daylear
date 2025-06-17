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

// GetRecipeRecipients retrieves the list of recipients for a recipe
func (s *RecipeService) GetRecipeRecipients(ctx context.Context, req *pb.GetRecipeRecipientsRequest) (*pb.RecipeRecipients, error) {
	user, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	mRecipe := model.Recipe{}
	_, err := s.recipeNamer.Parse(req.GetName(), &mRecipe)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", req.GetName())
	}

	mRecipe.Parent.UserId = user.Id.UserId

	recipients, err := s.domain.ListRecipeRecipients(ctx, mRecipe.Parent, mRecipe.Id)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to list recipe recipients")
		return nil, status.Errorf(codes.Internal, "unable to list recipe recipients")
	}

	// Convert recipients to proto format
	pbRecipients := make([]*pb.RecipeRecipients_Recipient, len(recipients))
	for i, recipient := range recipients {
		// Try to format as user first, then circle
		name, err := s.userNamer.Format(recipient.RecipeParent)
		if err != nil {
			name, err = s.circleNamer.Format(recipient.RecipeParent)
			if err != nil {
				s.log.Warn().Err(err).Msg("unable to format recipient name")
				return nil, status.Errorf(codes.Internal, "unable to format recipient name")
			}
		}

		pbRecipients[i] = &pb.RecipeRecipients_Recipient{
			Name:       name,
			Title:      recipient.Title,
			Permission: recipient.PermissionLevel,
		}
	}

	return &pb.RecipeRecipients{
		Name:       req.GetName(),
		Recipients: pbRecipients,
	}, nil
}
