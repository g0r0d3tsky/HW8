package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"net"
)

type Config struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

func (c *Config) ServerAddress() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Read() (*Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}
