package convert

import (
	model "github.com/jcfug8/daylear/server/core/model"
	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// ProtoToRecipe converts a protobuf Recipe to a model Recipe
func ProtoToRecipe(RecipeNamer namer.ReflectNamer[model.Recipe], proto *pb.Recipe) (model.Recipe, error) {
	recipe := model.Recipe{}
	var err error
	if proto.Name != "" {
		recipe, _, err = RecipeNamer.Parse(proto.Name, recipe)
		if err != nil {
			return recipe, err
		}
	}

	recipe.Title = proto.Title
	recipe.Description = proto.Description
	recipe.Directions = ProtosToDirections(proto.Directions)
	recipe.IngredientGroups = ProtosToIngredientGroups(proto.IngredientGroups)
	recipe.ImageURI = proto.ImageUri

	return recipe, nil
}

// RecipeToProto converts a model Recipe to a protobuf Recipe
func RecipeToProto(RecipeNamer namer.ReflectNamer[model.Recipe], recipe model.Recipe) (*pb.Recipe, error) {
	proto := &pb.Recipe{}
	name, err := RecipeNamer.Format(0, recipe)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	proto.Title = recipe.Title
	proto.Description = recipe.Description
	proto.Directions = DirectionsToProtos(recipe.Directions)
	proto.IngredientGroups = IngredientGroupsToProtos(recipe.IngredientGroups)
	proto.ImageUri = recipe.ImageURI

	return proto, nil
}

// RecipeListToProto converts a slice of model Recipes to a slice of protobuf OmniRecipes
func RecipeListToProto(RecipeNamer namer.ReflectNamer[model.Recipe], recipes []model.Recipe) ([]*pb.Recipe, error) {
	protos := make([]*pb.Recipe, len(recipes))
	for i, recipe := range recipes {
		proto := &pb.Recipe{}
		var err error
		if proto, err = RecipeToProto(RecipeNamer, recipe); err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToRecipe converts a slice of protobuf OmniRecipes to a slice of model Recipes
func ProtosToRecipe(RecipeNamer namer.ReflectNamer[model.Recipe], protos []*pb.Recipe) ([]model.Recipe, error) {
	res := make([]model.Recipe, len(protos))
	for i, proto := range protos {
		recipe := model.Recipe{}
		var err error
		if recipe, err = ProtoToRecipe(RecipeNamer, proto); err != nil {
			return nil, err
		}
		res[i] = recipe
	}
	return res, nil
}
