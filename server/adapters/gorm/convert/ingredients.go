package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

func IngredientFromCoreModel(m cmodel.Ingredient) gmodel.Ingredient {
	gm := gmodel.Ingredient{
		IngredientId: m.IngredientId,
		Title:        m.Title,
	}
	return gm
}

func IngredientToCoreModel(m gmodel.Ingredient) cmodel.Ingredient {
	cm := cmodel.Ingredient{
		IngredientId: m.IngredientId,
		Title:        m.Title,
	}
	return cm
}

func IngredientListFromCoreModel(m []cmodel.Ingredient) (res []gmodel.Ingredient) {
	res = make([]gmodel.Ingredient, len(m))
	for i, v := range m {
		res[i] = IngredientFromCoreModel(v)
	}
	return res
}

func IngredientListToCoreModel(m []gmodel.Ingredient) (res []cmodel.Ingredient) {
	res = make([]cmodel.Ingredient, len(m))
	for i, v := range m {
		res[i] = IngredientToCoreModel(v)
	}
	return res
}
