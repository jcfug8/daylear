package fieldmasker

import (
	"github.com/jcfug8/daylear/server/core/model"
)

var recipeFieldMap = map[string][]string{
	"name":                           {model.RecipeFields.Id},
	"title":                          {model.RecipeFields.Title},
	"description":                    {model.RecipeFields.Description},
	"directions":                     {model.RecipeFields.Directions},
	"ingredient_groups":              {model.RecipeFields.IngredientGroups},
	"image_uri":                      {model.RecipeFields.ImageURI},
	"visibility":                     {model.RecipeFields.VisibilityLevel},
	"recipe_access.name":             {model.RecipeFields.AccessId},
	"recipe_access.permission_level": {model.RecipeFields.PermissionLevel},
	"recipe_access.state":            {model.RecipeFields.State},
	"citation":                       {model.RecipeFields.Citation},
	"cook_duration":                  {model.RecipeFields.CookDurationSeconds},
	"prep_duration":                  {model.RecipeFields.PrepDurationSeconds},
	"total_duration":                 {model.RecipeFields.TotalDurationSeconds},
	"cooking_method":                 {model.RecipeFields.CookingMethod},
	"categories":                     {model.RecipeFields.Categories},
	"yield_amount":                   {model.RecipeFields.YieldAmount},
	"cuisines":                       {model.RecipeFields.Cuisines},
	"create_time":                    {model.RecipeFields.CreateTime},
	"update_time":                    {model.RecipeFields.UpdateTime},
}
