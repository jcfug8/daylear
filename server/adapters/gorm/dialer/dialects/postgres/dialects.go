package postgres

import (
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/rs/zerolog"
)

// NewDialects -
func NewDialects(driver Driver, log zerolog.Logger, configClient config.Client) (dialects []*Dialect) {
	// for _, c := range configs {
	dialects = append(dialects, NewDialect(DialectConfig{
		ConfigClient: configClient,
		Driver:       driver,
		Log:          log,
	}))
	// }

	return dialects
}
