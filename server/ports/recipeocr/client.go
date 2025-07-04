package recipeocr

import (
	"context"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/core/model"
)

type Client interface {
	OCRRecipe(ctx context.Context, files []file.File) (recipe model.Recipe, err error)
}
