package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// ListRecipes lists recipes.
func (repo *Client) ListRecipes(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Str("filter", filter).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("offset", offset).
		Logger()

	dbRecipes := []gmodel.Recipe{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "recipe.recipe_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.RecipeFieldMasker.Convert(fields)).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset))

	if authAccount.CircleId != 0 {
		tx = tx.Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_circle_id = ?", authAccount.CircleId).
			Joins("LEFT JOIN recipe_access as ra ON recipe.recipe_id = ra.recipe_id AND ra.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(recipe_access.recipient_circle_id = ? AND (recipe.visibility_level = ? OR ra.state = ?))",
				authAccount.CircleId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else if authAccount.UserId != 0 {
		tx = tx.Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN recipe_access as ra ON recipe.recipe_id = ra.recipe_id AND ra.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(recipe.visibility_level = ? OR (recipe_access.recipient_user_id = ? AND ra.state = ?))",
				types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else {
		tx = tx.Joins("LEFT JOIN recipe_access ON recipe.recipe_id = recipe_access.recipe_id AND recipe_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(recipe.visibility_level = ? OR recipe_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	conversion, err := gmodel.RecipeSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing recipe rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter"}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbRecipes).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list recipe rows")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Recipe, len(dbRecipes))
	for i, m := range dbRecipes {
		res[i], err = convert.RecipeToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("invalid recipe row when listing recipes")
			return nil, repository.ErrInternal{Msg: "invalid recipe row when listing recipes"}
		}
	}

	return res, nil
}

// CreateRecipe creates a new recipe.
func (repo *Client) CreateRecipe(ctx context.Context, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid recipe when creating recipe row")
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid recipe: %v", err)}
	}

	err = repo.db.
		Select(gmodel.RecipeFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create recipe row")
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid recipe row when creating recipe")
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}

// DeleteRecipe deletes a recipe.
func (repo *Client) DeleteRecipe(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.RecipeId) (cmodel.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", id.RecipeId).
		Logger()

	gm := gmodel.Recipe{RecipeId: id.RecipeId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.RecipeFieldMasker.GetAll()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete recipe row")
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err := convert.RecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid recipe row when deleting recipe")
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}

// GetRecipe gets a recipe.
func (repo *Client) GetRecipe(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.RecipeId, fields []string) (cmodel.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", id.RecipeId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.Recipe{}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.RecipeFieldMasker.Convert(fields)).
		Where("recipe.recipe_id = ?", id.RecipeId)

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get recipe row")
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err := convert.RecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid recipe row when getting recipe")
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}

// UpdateRecipe updates a recipe.
func (repo *Client) UpdateRecipe(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", m.Id.RecipeId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid recipe when updating recipe row")
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("error reading recipe: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.RecipeFieldMasker.Convert(fields)).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update recipe row")
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid recipe row when updating recipe")
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}
