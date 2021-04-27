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

// build global variables.
var Config Configuration

func (c *Configuration) Load() {
	config.Configure(c)
}

func (c *Configuration) Path() string {
	return config.FileNameFromFlag("apiconfig", "config.yaml", "file path to the config")
}

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
func init() {
	Config.Load()
}
