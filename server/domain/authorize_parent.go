package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/model"
)

// AuthorizeRecipeParent checks if the caller is authorized to access the parent's
// resources
func (d Domain) AuthorizeRecipeParent(ctx context.Context, tokenUser model.User, parent model.RecipeParent) error {
	if tokenUser.Id.UserId == parent.UserId {
		return nil
	}

	return fmt.Errorf("user not authorized")
}

// AuthorizeCircleParent checks if the caller is authorized to access the parent's
// resources
func (d Domain) AuthorizeCircleParent(ctx context.Context, tokenUser model.User, parent model.CircleParent) error {
	if tokenUser.Id.UserId == parent.UserId {
		return nil
	}

	return fmt.Errorf("user not authorized")
}
