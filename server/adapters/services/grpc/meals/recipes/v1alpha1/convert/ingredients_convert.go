package convert

import (
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// ProtoToRecipeIngredients converts a proto Ingredients to a domain Ingredients.
func ProtoToRecipeIngredients(proto *pb.Recipe_Ingredient) model.RecipeIngredient {
	ingredient := model.RecipeIngredient{}
	ingredient.Title = proto.Title
	ingredient.Optional = proto.Optional
	ingredient.MeasurementAmount = proto.MeasurementAmount
	ingredient.MeasurementType = proto.MeasurementType
	return ingredient
}

// RecipeIngredientsToProto converts a domain Ingredients to a proto Ingredients.
func RecipeIngredientsToProto(ingredient model.RecipeIngredient) *pb.Recipe_Ingredient {
	return &pb.Recipe_Ingredient{
		Title:             ingredient.Title,
		Optional:          ingredient.Optional,
		MeasurementAmount: ingredient.MeasurementAmount,
		MeasurementType:   ingredient.MeasurementType,
	}
}

// ProtosToRecipeIngredients converts a slice of proto Ingredients to a slice of domain Ingredients.
func ProtosToRecipeIngredients(protos []*pb.Recipe_Ingredient) []model.RecipeIngredient {
	ingredients := make([]model.RecipeIngredient, len(protos))
	for i, proto := range protos {
		ingredient := ProtoToRecipeIngredients(proto)
		ingredients[i] = ingredient
	}
	return ingredients
}

// RecipeIngredientsToProtos converts a slice of domain Ingredients to a slice of proto Ingredients.
func RecipeIngredientsToProtos(ingredients []model.RecipeIngredient) []*pb.Recipe_Ingredient {
	protos := make([]*pb.Recipe_Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		proto := RecipeIngredientsToProto(ingredient)
		protos[i] = proto
	}
	return protos
}

// ProtoToIngredientGroup converts a proto IngredientGroup to a domain IngredientGroup.
func ProtoToIngredientGroup(proto *pb.Recipe_IngredientGroup) model.IngredientGroup {
	ingredientGroup := model.IngredientGroup{}
	ingredientGroup.Title = proto.Title
	ingredientGroup.RecipeIngredients = ProtosToRecipeIngredients(proto.Ingredients)
	return ingredientGroup
}

// IngredientGroupToProto converts a domain IngredientGroup to a proto IngredientGroup.
func IngredientGroupToProto(ingredientGroup model.IngredientGroup) *pb.Recipe_IngredientGroup {
	return &pb.Recipe_IngredientGroup{
		Title:       ingredientGroup.Title,
		Ingredients: RecipeIngredientsToProtos(ingredientGroup.RecipeIngredients),
	}
}

// ProtosToIngredientGroup converts a slice of proto IngredientGroup to a slice of domain IngredientGroup.
func ProtosToIngredientGroups(protos []*pb.Recipe_IngredientGroup) []model.IngredientGroup {
	ingredientGroups := make([]model.IngredientGroup, len(protos))
	for i, proto := range protos {
		ingredientGroup := ProtoToIngredientGroup(proto)
		ingredientGroups[i] = ingredientGroup
	}
	return ingredientGroups
}

// IngredientGroupToProtos converts a slice of domain IngredientGroup to a slice of proto IngredientGroup.
func IngredientGroupsToProtos(ingredientGroups []model.IngredientGroup) []*pb.Recipe_IngredientGroup {
	protos := make([]*pb.Recipe_IngredientGroup, len(ingredientGroups))
	for i, ingredientGroup := range ingredientGroups {
		proto := IngredientGroupToProto(ingredientGroup)
		protos[i] = proto
	}
	return protos
}
