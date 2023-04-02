package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Addr        string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	StoragePath string `env:"FILE_STORAGE_PATH"`
}

type Flags struct {
	Addr        string
	BaseURL     string
	StoragePath string
}

func New() (*Config, error) {
	cfg := &Config{}

	flags := parseFlag()
	cfg.Addr = flags.Addr
	cfg.BaseURL = flags.BaseURL
	cfg.StoragePath = flags.StoragePath

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseFlag() *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.Addr, "a", "localhost:8080", "server address")
	flag.StringVar(&flags.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&flags.StoragePath, "f", "./storage.txt", "path to storage")
	flag.Parse()

	return flags
}
