package v1alpha1

import (
	"go.uber.org/fx"

	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

var Module = fx.Module(
	"recipeGrpcAdapter",
	fx.Provide(
		NewRecipeService,
		func(s *RecipeService) pb.RecipeServiceServer { return s },
		func(s *RecipeService) pb.RecipeAccessServiceServer { return s },
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name="v1alpha1RecipeAccessNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Recipe]() },
			fx.ResultTags(`name="v1alpha1RecipeNamer"`),
		),
	),
)
