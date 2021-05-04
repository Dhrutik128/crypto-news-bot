package main

import (
	"errors"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"time"
)

type Configuration struct {
	BotToken    string        `yaml:"bot_token" json:"bot_token"`
	RefreshRate time.Duration `yaml:"refresh_rate" json:"refresh_rate"`
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
	if c.RefreshRate == 0 {
		// default refresh rate is set to 1 hour
		c.RefreshRate = time.Hour
	} else {
		// refresh rate resolution is 1 minute
		c.RefreshRate = c.RefreshRate * time.Minute
	}
	if c.BotToken == "" {
		return errors.New("no telegram token provided")
	}
	return nil
}
