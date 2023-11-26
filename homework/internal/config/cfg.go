package config

import (
	"github.com/g0r0d3tsky/parser/pkg/cfgparser"
	"net"
)

func ServerAddress(cfg cfgparser.Config) string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

func Read(filePath string) (*cfgparser.Config, error) {
	config, err := cfgparser.ParseYAML(filePath)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
