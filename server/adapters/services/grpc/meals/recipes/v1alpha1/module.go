package v1alpha1

import (
	"go.uber.org/fx"

	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

var Module = fx.Module(
	"recipeGrpcAdapter",
	fx.Provide(
		NewRecipeService,
		func(s *RecipeService) pb.RecipeServiceServer { return s },
		func(s *RecipeService) pb.RecipeAccessServiceServer { return s },
	),
)
