package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// RecipeMap maps the core model fields to the database model fields for the unified Recipe model.
var RecipeMap = map[string]string{
	"permission": gmodel.RecipeAccessFields.PermissionLevel,
	"visibility": gmodel.RecipeFields.VisibilityLevel,
	"state":      gmodel.RecipeAccessFields.State,
}

// ListRecipes lists recipes.
func (repo *Client) ListRecipes(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string) ([]cmodel.Recipe, error) {
	dbRecipes := []gmodel.Recipe{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "recipe.recipe_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select("recipe.*, recipe_access.permission_level, recipe_access.state").
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset))

	if authAccount.CircleId != 0 {
		tx = tx.Where("(recipe_access.recipient_circle_id = ? OR recipe.visibility_level = ?) AND (recipe.visibility_level != ? OR recipe_access.permission_level = ?)",
			authAccount.CircleId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_circle_id = ?", authAccount.CircleId)
	} else {
		tx = tx.Where("(recipe_access.recipient_user_id = ? OR recipe.visibility_level = ?) AND (recipe.visibility_level != ? OR recipe_access.permission_level = ?)",
			authAccount.UserId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_user_id = ?", authAccount.UserId)
	}

	conversion, err := repo.recipeSQLConverter.Convert(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbRecipes).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Recipe, len(dbRecipes))
	for i, m := range dbRecipes {
		res[i], err = convert.RecipeToCoreModel(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read recipe: %v", err)
		}
	}

	return res, nil
}

// CreateRecipe creates a new recipe.
func (repo *Client) CreateRecipe(ctx context.Context, m cmodel.Recipe) (cmodel.Recipe, error) {
	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid recipe: %v", err)}
	}

	recipeFields := masks.RemovePaths(
		gmodel.RecipeFields.Mask(),
	)

	err = repo.db.
		Select(recipeFields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}

// DeleteRecipe deletes a recipe.
func (repo *Client) DeleteRecipe(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.RecipeId) (cmodel.Recipe, error) {
	gm := gmodel.Recipe{RecipeId: id.RecipeId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.RecipeFields.Mask()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err := convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}

// GetRecipe gets a recipe.
func (repo *Client) GetRecipe(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.RecipeId) (cmodel.Recipe, error) {
	gm := gmodel.Recipe{}

	tx := repo.db.WithContext(ctx).
		Select("recipe.*, recipe_access.permission_level").
		Where("recipe.recipe_id = ?", id.RecipeId)

	if authAccount.CircleId != 0 {
		tx = tx.Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_circle_id = ?", authAccount.CircleId).
			Where("(recipe.visibility_level = ? OR recipe_access.recipient_circle_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.CircleId)
	} else {
		tx = tx.Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_user_id = ?", authAccount.UserId).
			Where("(recipe.visibility_level = ? OR recipe_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId)
	}

	err := tx.First(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err := convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}

// UpdateRecipe updates a recipe.
func (repo *Client) UpdateRecipe(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("error reading recipe: %v", err)}
	}

	mask := masks.Map(fields, gmodel.RecipeMap)

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}
