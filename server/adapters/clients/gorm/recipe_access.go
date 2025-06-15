package gorm

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

func (repo *Client) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}

func (repo *Client) DeleteRecipeAccess(ctx context.Context, parent model.RecipeParent, id model.RecipeId) error {
	return nil
}

func (repo *Client) GetRecipeAccess(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}

func (repo *Client) ListRecipeAccesses(ctx context.Context, parent model.RecipeParent, pageSize int64, pageOffset int64, filter string) ([]model.RecipeAccess, error) {
	return nil, nil
}

func (repo *Client) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}
