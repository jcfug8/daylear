package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/model"
)

// AuthorizeRecipeParent checks if the caller is authorized to access the parent's
// resources
func (d Domain) AuthorizeRecipeParent(ctx context.Context, token string, parent model.RecipeParent) error {
	user, err := d.ParseToken(ctx, token)
	if err != nil {
		return err
	}

	if user.Id.UserId == parent.UserId {
		return nil
	}

	return fmt.Errorf("user not authorized")
}

// AuthorizeCircleParent checks if the caller is authorized to access the parent's
// resources
func (d Domain) AuthorizeCircleParent(ctx context.Context, token string, parent model.CircleParent) error {
	user, err := d.ParseToken(ctx, token)
	if err != nil {
		return err
	}

	if user.Id.UserId == parent.UserId {
		return nil
	}

	return fmt.Errorf("user not authorized")
}
