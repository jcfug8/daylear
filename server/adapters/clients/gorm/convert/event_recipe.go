package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// EventRecipeFromCoreModel converts a core EventRecipe to a GORM EventRecipe
func EventRecipeFromCoreModel(eventRecipe cmodel.EventRecipe) (gmodel.EventRecipe, error) {
	return gmodel.EventRecipe{
		EventRecipeId: eventRecipe.EventRecipeId.EventRecipeId,
		RecipeId:      eventRecipe.RecipeId.RecipeId,
		EventId:       eventRecipe.Parent.EventId,
		CreateTime:    eventRecipe.CreateTime,
	}, nil
}

// EventRecipeToCoreModel converts a GORM EventRecipe to a core EventRecipe
func EventRecipeToCoreModel(gormEventRecipe gmodel.EventRecipe) (cmodel.EventRecipe, error) {
	eventRecipe := cmodel.EventRecipe{
		EventRecipeId: cmodel.EventRecipeId{EventRecipeId: gormEventRecipe.EventRecipeId},
		RecipeId:      cmodel.RecipeId{RecipeId: gormEventRecipe.RecipeId},
		Parent: cmodel.EventRecipeParent{
			EventId: gormEventRecipe.EventId,
		},
		CreateTime: gormEventRecipe.CreateTime,
	}

	return eventRecipe, nil
}
