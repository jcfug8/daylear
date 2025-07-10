package fieldmasker

import (
	model "github.com/jcfug8/daylear/server/core/model"
)

var updateRecipeAccessFieldMap = map[string][]string{
	"level": {model.RecipeAccessFields.Level},
	"state": {model.RecipeAccessFields.State},
}
