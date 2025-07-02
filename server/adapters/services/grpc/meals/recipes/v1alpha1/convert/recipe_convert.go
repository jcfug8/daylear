package convert

import (
	model "github.com/jcfug8/daylear/server/core/model"
	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// ProtoToRecipe converts a protobuf Recipe to a model Recipe
func ProtoToRecipe(RecipeNamer namer.ReflectNamer, proto *pb.Recipe) (model.Recipe, error) {
	recipe := model.Recipe{}
	var err error
	if proto.Name != "" {
		_, err = RecipeNamer.Parse(proto.Name, &recipe)
		if err != nil {
			return recipe, err
		}
	}

	recipe.Title = proto.Title
	recipe.Description = proto.Description
	recipe.Directions = ProtosToDirections(proto.Directions)
	recipe.IngredientGroups = ProtosToIngredientGroups(proto.IngredientGroups)
	recipe.ImageURI = proto.ImageUri
	recipe.Visibility = proto.Visibility
	recipe.Permission = proto.Permission
	recipe.State = proto.State

	return recipe, nil
}

// RecipeToProto converts a model Recipe to a protobuf Recipe
func RecipeToProto(RecipeNamer namer.ReflectNamer, recipe model.Recipe) (*pb.Recipe, error) {
	proto := &pb.Recipe{}
	name, err := RecipeNamer.Format(recipe)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	proto.Title = recipe.Title
	proto.Description = recipe.Description
	proto.Directions = DirectionsToProtos(recipe.Directions)
	proto.IngredientGroups = IngredientGroupsToProtos(recipe.IngredientGroups)
	proto.ImageUri = recipe.ImageURI
	proto.Visibility = recipe.Visibility
	proto.Permission = recipe.Permission
	proto.State = recipe.State

	return proto, nil
}

// RecipeListToProto converts a slice of model Recipes to a slice of protobuf OmniRecipes
func RecipeListToProto(RecipeNamer namer.ReflectNamer, recipes []model.Recipe) ([]*pb.Recipe, error) {
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
func ProtosToRecipe(RecipeNamer namer.ReflectNamer, protos []*pb.Recipe) ([]model.Recipe, error) {
	res := make([]model.Recipe, len(protos))
	for i, proto := range protos {
		var err error
		res[i], err = ProtoToRecipe(RecipeNamer, proto)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
