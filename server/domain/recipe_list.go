package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
)

// ListRecipes lists recipes.
func (d *Domain) ListRecipes(ctx context.Context, page *model.PageToken[model.Recipe], parent model.RecipeParent, filter string, fieldMask []string) (recipes []model.Recipe, err error) {
	filter = fmt.Sprintf("%s user_id = %d", filter, parent.UserId)

	recipes, err = d.repo.ListRecipes(ctx, page, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	for i := range recipes {
		recipes[i].Parent = parent
	}

	return recipes, nil
}
