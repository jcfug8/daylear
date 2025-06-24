package gorm

import (
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/ports/repository"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var _ repository.Client = (*Client)(nil)
var _ repository.TxClient = (*Client)(nil)

// ClientParams defines the dependencies for the GORM client.
type ClientParams struct {
	fx.In

	DB  *gorm.DB
	Log zerolog.Logger
}

// NewClient creates a new GORM client.
func NewClient(p ClientParams) (*Client, error) {
	return &Client{
		db:                             p.DB,
		log:                            p.Log,
		recipeAccessUserSQLConverter:   filter.NewSQLConverter(RecipeAccessUserMap, true),
		recipeAccessCircleSQLConverter: filter.NewSQLConverter(RecipeAccessCircleMap, true),
	}, nil
}

// Client defines a GORM client.
type Client struct {
	db    *gorm.DB
	level int
	log   zerolog.Logger

	recipeAccessUserSQLConverter   *filter.SQLConverter
	recipeAccessCircleSQLConverter *filter.SQLConverter
}

// Migrate migrates the database.
func (repo *Client) Migrate() (err error) {
	for _, m := range model.AllModels() {
		repo.log.Info().Msgf("auto migrating model %T", m)
		if err = repo.db.AutoMigrate(m); err != nil {
			return ConvertGormError(err)
		}
	}

	return nil
}
