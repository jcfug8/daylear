package convert

import (
	model "github.com/jcfug8/daylear/server/core/model"
	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToRecipe converts a protobuf Recipe to a model Recipe
func ProtoToRecipe(RecipeNamer namer.ReflectNamer, proto *pb.Recipe) (int, model.Recipe, error) {
	recipe := model.Recipe{}
	var nameIndex int
	var err error
	if proto.Name != "" {
		nameIndex, err = RecipeNamer.Parse(proto.Name, &recipe)
		if err != nil {
			return nameIndex, recipe, err
		}
	}

	recipe.Title = proto.Title
	recipe.Description = proto.Description
	recipe.Directions = ProtosToDirections(proto.Directions)
	recipe.IngredientGroups = ProtosToIngredientGroups(proto.IngredientGroups)
	recipe.ImageURI = proto.ImageUri
	recipe.VisibilityLevel = proto.Visibility
	recipe.Citation = proto.Citation
	if proto.CookDuration != nil {
		recipe.CookDuration = proto.CookDuration.AsDuration()
	}
	if proto.PrepDuration != nil {
		recipe.PrepDuration = proto.PrepDuration.AsDuration()
	}
	if proto.TotalDuration != nil {
		recipe.TotalDuration = proto.TotalDuration.AsDuration()
	}
	recipe.CookingMethod = proto.CookingMethod
	recipe.Categories = proto.Categories
	recipe.YieldAmount = proto.YieldAmount
	recipe.Cuisines = proto.Cuisines
	if proto.CreateTime != nil {
		recipe.CreateTime = proto.CreateTime.AsTime()
	}
	if proto.UpdateTime != nil {
		recipe.UpdateTime = proto.UpdateTime.AsTime()
	}

	return nameIndex, recipe, nil
}

// RecipeToProto converts a model Recipe to a protobuf Recipe
func RecipeToProto(RecipeNamer namer.ReflectNamer, AccessNamer namer.ReflectNamer, recipe model.Recipe, options ...namer.FormatReflectNamerOption) (*pb.Recipe, error) {
	proto := &pb.Recipe{}

	if recipe.Id.RecipeId != 0 {
		name, err := RecipeNamer.Format(recipe, options...)
		if err != nil {
			return proto, err
		}
		proto.Name = name
	}

	proto.Title = recipe.Title
	proto.Description = recipe.Description
	proto.Directions = DirectionsToProtos(recipe.Directions)
	proto.IngredientGroups = IngredientGroupsToProtos(recipe.IngredientGroups)
	proto.ImageUri = recipe.ImageURI
	proto.Visibility = recipe.VisibilityLevel
	proto.Citation = recipe.Citation
	proto.CookDuration = durationpb.New(recipe.CookDuration)
	proto.PrepDuration = durationpb.New(recipe.PrepDuration)
	proto.TotalDuration = durationpb.New(recipe.TotalDuration)
	proto.CookingMethod = recipe.CookingMethod
	proto.Categories = recipe.Categories
	proto.YieldAmount = recipe.YieldAmount
	proto.Cuisines = recipe.Cuisines
	proto.CreateTime = timestamppb.New(recipe.CreateTime)
	proto.UpdateTime = timestamppb.New(recipe.UpdateTime)

	// Handle recipe_access field if present
	if (recipe.RecipeAccess != model.RecipeAccess{}) {
		name, err := AccessNamer.Format(recipe.RecipeAccess)
		if err == nil {
			proto.RecipeAccess = &pb.Recipe_RecipeAccess{
				Name:            name,
				PermissionLevel: recipe.RecipeAccess.PermissionLevel,
				State:           recipe.RecipeAccess.State,
				AcceptTarget:    recipe.RecipeAccess.AcceptTarget,
			}
		}
	}

	return proto, nil
}

// RecipeListToProto converts a slice of model Recipes to a slice of protobuf OmniRecipes
func RecipeListToProto(RecipeNamer namer.ReflectNamer, AccessNamer namer.ReflectNamer, recipes []model.Recipe) ([]*pb.Recipe, error) {
	protos := make([]*pb.Recipe, len(recipes))
	for i, recipe := range recipes {
		proto := &pb.Recipe{}
		var err error
		if proto, err = RecipeToProto(RecipeNamer, AccessNamer, recipe); err != nil {
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
		_, res[i], err = ProtoToRecipe(RecipeNamer, proto)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
