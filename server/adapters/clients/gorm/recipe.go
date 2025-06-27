package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// ListRecipes lists recipes.
func (repo *Client) ListRecipes(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.Recipe, error) {
	queryModel := gmodel.Recipe{}

	args := make([]any, 0, 1)

	fields = masks.Map(fields, gmodel.RecipeMap)

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		for i, field := range fields {
			fields[i] = fmt.Sprintf("r.%s", field)
		}
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("r.recipe_id", "="),
		})

	filterClause, _, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if authAccount.CircleId != 0 {
		tx.Joins("JOIN recipe_circle ON recipe_circle.recipe_id = r.recipe_id AND recipe_circle.circle_id = ?", authAccount.CircleId)
	} else if authAccount.UserId != 0 {
		tx.Joins("JOIN recipe_user ON recipe_user.recipe_id = r.recipe_id AND recipe_user.user_id = ?", authAccount.UserId)
	}

	if authAccount.UserId != 0 && authAccount.CircleId != 0 {
		tx.Joins("JOIN circle_user ON circle_user.circle_id = recipe_circle.circle_id AND circle_user.user_id = ?", authAccount.UserId)
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}
	if len(args) > 0 {
		tx = tx.Where(queryModel, args...)
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "r.recipe_id"},
		Desc:   true,
	}}

	tx = tx.Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset))

	var mods []gmodel.Recipe
	tableAlias := fmt.Sprintf("%s AS r", gmodel.Recipe{}.TableName())
	if err = tx.Table(tableAlias).Find(&mods).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Recipe, len(mods))
	for i, m := range mods {
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
func (repo *Client) DeleteRecipe(ctx context.Context, id cmodel.RecipeId) (cmodel.Recipe, error) {
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
func (repo *Client) GetRecipe(ctx context.Context, id cmodel.RecipeId, fields []string) (cmodel.Recipe, error) {
	gm := gmodel.Recipe{RecipeId: id.RecipeId}

	mask := masks.Map(fields, gmodel.RecipeMap)
	if len(mask) == 0 {
		mask = gmodel.RecipeFields.Mask()
	}

	err := repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		First(&gm).Error
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
func (repo *Client) UpdateRecipe(ctx context.Context, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
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
