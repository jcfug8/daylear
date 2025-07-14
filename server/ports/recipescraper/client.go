package recipescraper

import (
	"context"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/core/model"
)

type Client interface {
	RecipeFromData(ctx context.Context, data []byte) (model.Recipe, error)
	RecipeFromImage(ctx context.Context, files []file.File) (recipe model.Recipe, err error)
}
