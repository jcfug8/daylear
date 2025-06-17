package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

func (d *Domain) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	return model.RecipeAccess{}, nil
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, parent model.RecipeId, id model.RecipeAccessId) error {
	return nil
}

func (d *Domain) GetRecipeAccess(ctx context.Context, parent model.RecipeId, id model.RecipeAccessId) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, parent model.RecipeId, pageSize int32, pageOffset int32, filter string) ([]model.RecipeAccess, error) {
	return nil, nil
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}
