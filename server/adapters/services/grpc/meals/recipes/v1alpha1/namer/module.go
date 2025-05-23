package namer

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"recipe_namer",
	fx.Provide(
		func() (namer.ReflectNamer[model.Recipe], error) {
			return namer.NewReflectNamer[model.Recipe](
				&pb.Recipe{},
			)
		},
	),
)
