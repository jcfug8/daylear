package gorm

import (
	"github.com/jcfug8/daylear/server/ports/repository"

	dialer "github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer"

	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"repository",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(repository.Client)),
			fx.OnStart(func(client repository.Client) error {
				return client.Migrate()
			}),
		),
	),

	// IRIOMO:CUSTOM_CODE_SLOT_START repositoryImports
	dialer.ProvideConfigOptions(dialer.DisableForeignKeyConstraintWhenMigrating()),
	// IRIOMO:CUSTOM_CODE_SLOT_END

	dialer.ProvideConfigOptions(dialer.TranslateError()),
)
