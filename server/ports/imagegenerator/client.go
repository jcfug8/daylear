package imagegenerator

import (
	"context"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/core/model"
)

type Client interface {
	GenerateRecipeImage(ctx context.Context, recipe model.Recipe) (file.File, error)
}
