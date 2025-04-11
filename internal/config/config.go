package config

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	BotToken string `long:"token" env:"BOT_TOKEN" required:"true" description:"telegram bot token"`
}

func Load() (Config, error) {
	var config Config
	p := flags.NewParser(&config, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		return Config{}, fmt.Errorf("parse: %v", err)
	}

	return config, nil
}
