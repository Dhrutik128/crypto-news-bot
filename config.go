package main

import (
	"errors"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"time"
)

type Configuration struct {
	BotToken              string        `yaml:"bot_token" json:"bot_token"`
	RefreshPeriodDuration time.Duration `yaml:"refresh_period_duration" json:"refresh_period_duration"`
	NewsStorageDuration   time.Duration `yaml:"news_storage_duration" json:"news_storage_duration"`
}

// Config is the global bot configuration read from file.
var Config Configuration

func init() {
	Config.Load()
}

// Load the configuration
func (c *Configuration) Load() {
	config.Configure(c)
}

// Path must return the path to configuration file (yaml or json)
func (c *Configuration) Path() string {
	return config.FileNameFromFlag("botconfig", "config.yaml", "file path to the config")
}

// Valid checks if configuration is valid
func (c *Configuration) Valid() error {
	if c.RefreshPeriodDuration == 0 {
		// default refresh rate is set to 1 hour
		c.RefreshPeriodDuration = time.Hour
	} else {
		// refresh rate resolution is 1 minute
		c.RefreshPeriodDuration = c.RefreshPeriodDuration * time.Minute
	}
	if c.BotToken == "" {
		return errors.New("no telegram token provided")
	}
	return nil
}
