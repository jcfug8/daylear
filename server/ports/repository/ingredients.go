package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the ingredient in the database.
type ingredientClient interface {
	BulkCreateIngredients(context.Context, []model.Ingredient) ([]model.Ingredient, error)
	BulkDeleteIngredients(context.Context, string) ([]model.Ingredient, error)
}
