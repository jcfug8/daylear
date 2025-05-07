package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// ListRecipeUsers lists recipes.
func (repo *Client) BulkDeleteRecipeUsers(ctx context.Context, filter string) error {
	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("recipe_id", "="),
			"user_id":   filtering.NewSQLField[int64]("user_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	var dbRecipeUsers []gmodel.RecipeUser
	if err = tx.Delete(&dbRecipeUsers).Error; err != nil {
		return repository.ErrInternal{Msg: fmt.Sprintf("failed to delete recipe users: %v", err)}
	}

	return nil
}
