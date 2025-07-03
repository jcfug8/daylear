package recipescraper

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type HostSpecificClient interface {
	GetHost() []string
	DefaultClient
}

type DefaultClient interface {
	ScrapeRecipe(ctx context.Context, uri string) (model.Recipe, error)
}
