package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/migrations"
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
		db:                       p.DB,
		log:                      p.Log,
		recipeAccessSQLConverter: filter.NewSQLConverter(RecipeAccessMap, true),
		recipeSQLConverter:       filter.NewSQLConverter(RecipeMap, true),
		circleSQLConverter:       filter.NewSQLConverter(CircleMap, true),
		circleAccessSQLConverter: filter.NewSQLConverter(CircleAccessMap, true),
		userSQLConverter:         filter.NewSQLConverter(UserMap, true),
		userCircleSQLConverter:   filter.NewSQLConverter(UserCircleMap, true),
	}, nil
}

// Client defines a GORM client.
type Client struct {
	db    *gorm.DB
	level int
	log   zerolog.Logger

	recipeAccessSQLConverter *filter.SQLConverter
	recipeSQLConverter       *filter.SQLConverter
	circleSQLConverter       *filter.SQLConverter
	circleAccessSQLConverter *filter.SQLConverter
	userSQLConverter         *filter.SQLConverter
	userCircleSQLConverter   *filter.SQLConverter
}

// Migrate migrates the database.
func (repo *Client) Migrate() (err error) {
	repo.log.Info().Msg("migrating database")

	err = migrations.RunMigrations(context.Background(), repo.db, manualMigrations)
	if err != nil {
		repo.log.Error().Err(err).Msg("failed to migrate database")
		return errors.New("failed to migrate database")
	}

	for _, m := range model.AllModels() {
		repo.log.Info().Msgf("auto migrating model %T", m)
		if err = repo.db.AutoMigrate(m); err != nil {
			return ConvertGormError(err)
		}
	}

	return nil
}
