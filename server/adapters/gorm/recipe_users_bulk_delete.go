package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"gorm.io/gorm/clause"
)

// ListRecipeUsers lists recipes.
func (repo *Client) BulkDeleteRecipeUsers(ctx context.Context, filter string) error {
	errz := errz.Context("repository.list_recipes")

	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("recipe_id", "="),
			"user_id":   filtering.NewSQLField[int64]("user_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return errz.NewInvalidArgument("invalid filter: %v", err)
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	var dbRecipeUsers []gmodel.RecipeUser
	if err = tx.Delete(&dbRecipeUsers).Error; err != nil {
		return ErrzError(errz, "", err)
	}

	return nil
}
