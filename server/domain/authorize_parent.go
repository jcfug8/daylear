package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/model"
)

// AuthorizeParent checks if the caller is authorized to access the parent's
// resources
func (d Domain) AuthorizeParent(ctx context.Context, token string, parent model.RecipeParent) error {
	user, err := d.ParseToken(ctx, token)
	if err != nil {
		return err
	}

	if user.Id.UserId == parent.UserId {
		return nil
	}

	return fmt.Errorf("user not authorized")
}
