package logger

import "github.com/rs/zerolog"

// Plugin -
type Plugin interface {
	Logger(zerolog.Logger) zerolog.Logger
}

// ApplyPlugins -
func ApplyPlugins(log zerolog.Logger, plugins []Plugin) zerolog.Logger {
	for _, plugin := range plugins {
		log = plugin.Logger(log)
	}

	return log
}
