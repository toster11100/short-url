package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Addr    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func New() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
