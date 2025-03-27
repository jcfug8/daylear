package dialer

import (
	"time"

	"github.com/jcfug8/daylear/server/adapters/gorm/dialer/logger"
	"github.com/jcfug8/daylear/server/adapters/gorm/dialer/plugins"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ConfigParams -
type ConfigParams struct {
	fx.In

	Log     zerolog.Logger
	Opts    []ConfigOption `group:"gormDialerConfigOptions"`
	Plugins []gorm.Plugin  `group:"gormDialerPlugins"`
}

// NewConfig -
func NewConfig(p ConfigParams) *Config {
	config := &Config{
		Config: &gorm.Config{
			Logger: logger.New(p.Log, nil),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
		log: p.Log,
		plugins: append([]gorm.Plugin{
			&plugins.PGCrypto{},
			plugins.NewSnowflake(p.Log),
		}, p.Plugins...),
	}

	for _, opt := range p.Opts {
		opt.apply(config.Config)
	}

	return config
}

// Config -
type Config struct {
	*gorm.Config
	log     zerolog.Logger
	plugins []gorm.Plugin
}

// ---------------------------------------------------------------

// ConfigOption -
type ConfigOption interface {
	apply(*gorm.Config)
}

func newFuncConfigOption(f func(*gorm.Config)) ConfigOption {
	return &funcConfigOption{f: f}
}

type funcConfigOption struct {
	f func(*gorm.Config)
}

func (fco *funcConfigOption) apply(c *gorm.Config) {
	fco.f(c)
}

// ---------------------------------------------------------------

// SkipDefaultTransaction -
func SkipDefaultTransaction() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.SkipDefaultTransaction = true
	})
}

// WithNamingStrategy -
func WithNamingStrategy(namer schema.Namer) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.NamingStrategy = namer
	})
}

// FullSaveAssociations -
func FullSaveAssociations() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.FullSaveAssociations = true
	})
}

// WithLogger -
func WithLogger(log zerolog.Logger, config *glogger.Config) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.Logger = logger.New(log, config)
	})
}

// WithNowFunc -
func WithNowFunc(f func() time.Time) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.NowFunc = f
	})
}

// DryRun -
func DryRun() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.DryRun = true
	})
}

// PrepareStmt -
func PrepareStmt() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.PrepareStmt = true
	})
}

// DisableAutomaticPing -
func DisableAutomaticPing() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.DisableAutomaticPing = true
	})
}

// DisableForeignKeyConstraintWhenMigrating -
func DisableForeignKeyConstraintWhenMigrating() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.DisableForeignKeyConstraintWhenMigrating = true

		// Adding callback to remove on conflict clauses.
		// Without doing this if foreign keys are disabled gorm will fail to create an object with associations.
		// This is because gorm generates invalid sql - the generated on conflict clause does not specify a conflict to check for
		removeOnConflictPlugin := new(plugins.RemoveOnConflict)
		// Making it nil safe
		if c.Plugins == nil {
			c.Plugins = map[string]gorm.Plugin{}
		}
		c.Plugins[removeOnConflictPlugin.Name()] = removeOnConflictPlugin
	})
}

// IgnoreRelationshipsWhenMigrating -
func IgnoreRelationshipsWhenMigrating() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.IgnoreRelationshipsWhenMigrating = true
	})
}

// DisableNestedTransaction -
func DisableNestedTransaction() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.DisableNestedTransaction = true
	})
}

// AllowGlobalUpdate -
func AllowGlobalUpdate() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.AllowGlobalUpdate = true
	})
}

// QueryFields -
func QueryFields() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.QueryFields = true
	})
}

// CreateBatchSize -
func CreateBatchSize(size int) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.CreateBatchSize = size
	})
}

// TranslateError -
func TranslateError() ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.TranslateError = true
	})
}

// WithClauseBuilders -
func WithClauseBuilders(builders map[string]clause.ClauseBuilder) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.ClauseBuilders = builders
	})
}

// WithConnPool -
func WithConnPool(pool gorm.ConnPool) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.ConnPool = pool
	})
}

// WithPlugins - Note that this will replace/overwrie any plugins which were previously added
func WithPlugins(plugins map[string]gorm.Plugin) ConfigOption {
	return newFuncConfigOption(func(c *gorm.Config) {
		c.Plugins = plugins
	})
}
