package config

import (
	"os"

	"sidecar/route"
	"sidecar/store"
)

type Config struct {
	Routes  []route.Route     `json:"routes" yaml:"routes"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	APIKey  route.APIKey      `json:"api_key" yaml:"api_key"`
	Debug   bool              `json:"debug" yaml:"debug"`
	Verbose bool              `json:"verbose" yaml:"verbose"`
}

func Find() (Config, error) {
	p, err := store.Find(store.Config)
	if err != nil {
		return Config{}, err
	}
	var c Config
	err = store.Load(p, &c)
	if err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}
	return c, nil
}

func Load(p string) (Config, error) {
	var c Config
	err := store.Load(p, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func Store(p string, c Config) error {
	return store.Store(p, c)
}
