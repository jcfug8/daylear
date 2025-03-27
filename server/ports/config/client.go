package config

// Client -
type Client interface {
	GetConfig() map[string]any
}
