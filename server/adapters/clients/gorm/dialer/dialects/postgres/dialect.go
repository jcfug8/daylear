package postgres

import (
	"fmt"
	"time"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/rs/zerolog"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ dialer.Dialect = (*Dialect)(nil)

// DialectConfig -
type DialectConfig struct {
	fx.In

	ConfigClient config.Client
	Driver       Driver `optional:"true"`
	Log          zerolog.Logger
}

// NewDialect -
func NewDialect(cfg DialectConfig) *Dialect {
	var err error
	config := cfg.ConfigClient.GetConfig()["postgres"].(map[string]interface{})

	host := config["host"].(string)
	port := config["port"].(string)
	user := config["user"].(string)
	password := config["password"].(string)
	dbname := config["dbname"].(string)
	sslmode := config["sslmode"].(string)

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)
	// maxIdleConns := int(config["maxIdleConns"].(float64))
	// maxOpenConns := int(config["maxOpenConns"].(float64))
	var maxConnLifetime time.Duration
	connLifeTime, ok := config["maxConnLifetime"]
	if ok {
		maxConnLifetime, err = time.ParseDuration(connLifeTime.(string))
		if err != nil {
			cfg.Log.Panic().Msgf("error parsing maxConnLifetime: %s", err)
		}
	}
	var maxConnIdleTime time.Duration
	connIdleTime, ok := config["maxConnIdleTime"]
	if ok {
		maxConnIdleTime, err = time.ParseDuration(connIdleTime.(string))
		if err != nil {
			cfg.Log.Panic().Msgf("error parsing maxConnIdleTime: %s", err)
		}
	}

	return &Dialect{
		dns: dns,
		// maxIdleConns:    maxIdleConns,
		// maxOpenConns:    maxOpenConns,
		maxConnLifetime: maxConnLifetime,
		maxConnIdleTime: maxConnIdleTime,
		driver:          cfg.Driver,
		log:             cfg.Log,
	}
}

// Dialect -
type Dialect struct {
	dns             string
	maxIdleConns    int
	maxOpenConns    int
	maxConnLifetime time.Duration
	maxConnIdleTime time.Duration
	driver          Driver
	log             zerolog.Logger
}

// Dialector implements Dialect
func (d *Dialect) Dialector() gorm.Dialector {
	config := postgres.Config{
		DSN: d.dns,
	}

	if d.driver == nil {
		return postgres.New(config)
	}

	return &pgDriver{
		Dialector: postgres.New(config).(*postgres.Dialector),
		driver:    d.driver,
	}
}

// Options implements Dialect
func (d *Dialect) Options() []gorm.Option {
	d.log.Info().Msgf("creating postgres options MaxIdleConns: %d , MaxOpenConns: %d , MaxConnLifetime: %d, MaxConnIdleTime %d", d.maxIdleConns, d.maxOpenConns, d.maxConnLifetime, d.maxConnIdleTime)
	return []gorm.Option{
		dialer.NewPoolConfigOption(
			d.maxIdleConns,
			d.maxOpenConns,
			d.maxConnLifetime,
			d.maxConnIdleTime,
		),
	}
}
