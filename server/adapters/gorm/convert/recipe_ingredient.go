package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

func RecipeIngredientFromCoreModel(m cmodel.RecipeIngredient) gmodel.RecipeIngredient {
	gm := gmodel.RecipeIngredient{
		RecipeIngredientId:   m.RecipeIngredientId,
		IngredientId:         m.IngredientId,
		MeasurementAmount:    m.MeasurementAmount,
		MeasurementType:      m.MeasurementType,
		IngredientGroupIndex: m.IngredientGroupIndex,
	}
	return gm
}

func RecipeIngredientToCoreModel(m gmodel.RecipeIngredient) cmodel.RecipeIngredient {
	cm := cmodel.RecipeIngredient{
		RecipeIngredientId:   m.RecipeIngredientId,
		MeasurementAmount:    m.MeasurementAmount,
		MeasurementType:      m.MeasurementType,
		IngredientGroupIndex: m.IngredientGroupIndex,
		Ingredient: cmodel.Ingredient{
			IngredientId: m.IngredientId,
		},
	}
	return cm
}

func RecipeIngredientListFromCoreModel(recipeId cmodel.RecipeId, m []cmodel.RecipeIngredient) (res []gmodel.RecipeIngredient) {
	res = make([]gmodel.RecipeIngredient, len(m))
	for i, v := range m {
		res[i] = RecipeIngredientFromCoreModel(v)
		res[i].RecipeId = recipeId.RecipeId
	}
	return res
}

func RecipeIngredientListToCoreModel(m []gmodel.RecipeIngredient) (res []cmodel.RecipeIngredient) {
	res = make([]cmodel.RecipeIngredient, len(m))
	for i, v := range m {
		res[i] = RecipeIngredientToCoreModel(v)
	}
	return res
}

func IngredientGroupsFromCoreModel(recipeId cmodel.RecipeId, m []cmodel.IngredientGroup) []gmodel.RecipeIngredient {
	res := []gmodel.RecipeIngredient{}
	for _, group := range m {
		recipeIngredients := RecipeIngredientListFromCoreModel(recipeId, group.RecipeIngredients)
		res = append(res, recipeIngredients...)
	}
	return res
}

func IngredientGroupsToCoreModel(m []gmodel.RecipeIngredient) []cmodel.IngredientGroup {
	res := make([]cmodel.IngredientGroup, len(m))
	for i, v := range m {
		res[i] = cmodel.IngredientGroup{
			RecipeIngredients: []cmodel.RecipeIngredient{RecipeIngredientToCoreModel(v)},
		}
	}
	return res
}
