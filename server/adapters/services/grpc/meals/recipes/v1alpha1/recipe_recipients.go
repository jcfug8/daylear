package v1alpha1

import (
	"context"

	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// GetRecipeRecipients retrieves the list of recipients for a recipe
func (s *RecipeService) GetRecipeRecipients(ctx context.Context, req *pb.GetRecipeRecipientsRequest) (*pb.RecipeRecipients, error) {
	return nil, nil
}
