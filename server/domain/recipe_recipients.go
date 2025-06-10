package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// ListRecipeRecipients retrieves all users and circles with access to a recipe.
// It requires a valid recipe ID and parent user ID.
func (d *Domain) ListRecipeRecipients(ctx context.Context, parent model.RecipeParent, id model.RecipeId) ([]model.RecipeRecipient, error) {
	if id.RecipeId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	if parent.UserId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	_, err := d.repo.GetRecipeRecipient(ctx, parent, id)
	if err != nil {
		return nil, err
	}
	if parent.CircleId != 0 {
		_, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
		if err != nil {
			return nil, err
		}
	}

	// get recipe recipients
	recipients, err := d.repo.ListRecipeRecipients(ctx, id)
	if err != nil {
		return nil, err
	}

	return recipients, nil
}
