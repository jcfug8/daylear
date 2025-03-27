package convert

import (
	"encoding/json"

	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// RecipeFromCoreModel converts a core model to a gorm model.
func RecipeFromCoreModel(m cmodel.Recipe) (gmodel.Recipe, error) {
	var err error
	recipe := gmodel.Recipe{
		RecipeId:    m.Id.RecipeId,
		Title:       m.Title,
		Description: m.Description,
		ImageURI:    m.ImageURI,
	}

	recipe.Directions, err = json.Marshal(m.Directions)
	if err != nil {
		return gmodel.Recipe{}, err
	}

	recipe.IngredientGroups, err = json.Marshal(m.IngredientGroups)
	if err != nil {
		return gmodel.Recipe{}, err
	}

	return recipe, nil
}

// RecipeToCoreModel converts a gorm model to a core model.
func RecipeToCoreModel(m gmodel.Recipe) (cmodel.Recipe, error) {
	var err error
	recipe := cmodel.Recipe{
		Id: cmodel.RecipeId{
			RecipeId: m.RecipeId,
		},
		Title:       m.Title,
		Description: m.Description,
		ImageURI:    m.ImageURI,
	}

	if m.Directions != nil {
		err = json.Unmarshal(m.Directions, &recipe.Directions)
		if err != nil {
			return cmodel.Recipe{}, err
		}
	}

	if m.IngredientGroups != nil {
		err = json.Unmarshal(m.IngredientGroups, &recipe.IngredientGroups)
		if err != nil {
			return cmodel.Recipe{}, err
		}
	}

	return recipe, nil
}

// RecipeListFromCoreModel converts a list of core models to a list of gorm models.
func RecipeListFromCoreModel(m []cmodel.Recipe) (res []gmodel.Recipe, err error) {
	res = make([]gmodel.Recipe, len(m))
	for i, v := range m {
		res[i], err = RecipeFromCoreModel(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// RecipeListToCoreModel converts a list of gorm models to a list of core models.
func RecipeListToCoreModel(m []gmodel.Recipe) (res []cmodel.Recipe, err error) {
	res = make([]cmodel.Recipe, len(m))
	for i, v := range m {
		res[i], err = RecipeToCoreModel(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
