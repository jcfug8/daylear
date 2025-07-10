package ingredientcleaner

import "context"

type Client interface {
	CleanIngredients(ctx context.Context, ingredients []string) (cleanedIngredients []string, err error)
}
