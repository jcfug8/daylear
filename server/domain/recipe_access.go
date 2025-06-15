package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

func (d *Domain) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, parent model.RecipeParent, id model.RecipeId) error {
	return nil
}

func (d *Domain) GetRecipeAccess(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, parent model.RecipeParent, pageSize int64, pageOffset int64, filter string) ([]model.RecipeAccess, error) {
	return nil, nil
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}
