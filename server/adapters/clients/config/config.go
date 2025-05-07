package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var (
	configFilename = "daylear.yaml"
	envFilename    = ".env"
)

const (
	// EnvPrefix -
	EnvPrefix = "DAYLEAR"
	PORT      = "PORT"
)

type Client struct {
	config map[string]interface{}
	log    zerolog.Logger
}

// NewConfig -
func NewConfig(log zerolog.Logger) Client {
	client := Client{
		log:    log,
		config: make(map[string]interface{}),
	}

	client.loadFile()
	client.loadEnv()

	return client
}

// GetConfig -
func (c Client) GetConfig() map[string]interface{} {
	return c.config
}

func (c *Client) loadEnv() {
	envFilePath, err := c.findFile(envFilename)
	if err != nil {
		c.log.Info().Msg("no env file found")
	} else {
		// Load environment variables from .env file
		if err := godotenv.Load(envFilePath); err != nil {
			c.log.Info().Msg("No .env file found")
		}
	}

	// set all other env variables
	for _, keyValue := range os.Environ() {
		splitKeyValue := strings.Split(keyValue, "=")
		key := splitKeyValue[0]
		value := splitKeyValue[1]
		if !strings.HasPrefix(key, EnvPrefix) {
			continue
		}
		key = strings.TrimPrefix(key, EnvPrefix+"_")
		keyParts := strings.Split(strings.ToLower(key), "_")
		c.set(c.config, keyParts, value)
	}
}

func (c *Client) loadFile() {
	if c.config == nil {
		c.config = make(map[string]interface{})
	}

	configFilePath, err := c.findFile(configFilename)
	if err != nil {
		c.log.Info().Msg("No config file found")
		return
	}

	c.log.Info().Msgf("loading config from %s", configFilePath)

	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		c.log.Panic().Err(err).Msg("failed to read config file")
	}

	err = yaml.Unmarshal(configFileBytes, &c.config)
	if err != nil {
		c.log.Panic().Err(err).Msg("failed to unmarshal config file")
	}
}

func (c *Client) set(config map[string]interface{}, keyParts []string, value string) {
	if len(keyParts) == 1 {
		config[keyParts[0]] = value
		return
	}

	subConfig, ok := c.config[keyParts[0]]
	if !ok {
		subConfig = make(map[string]interface{})
		config[keyParts[0]] = subConfig
	}

	c.set(subConfig.(map[string]interface{}), keyParts[1:], value)
}

func (c *Client) findFile(filename string) (string, error) {
	var foundFilePath string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			foundFilePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if foundFilePath == "" {
		return "", fmt.Errorf("file not found")
	}

	return foundFilePath, nil
}
