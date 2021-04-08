package config

import (
	"os"

	"sidecar/flags"
	"sidecar/route"
	"sidecar/store"
)

type Config struct {
	Routes  []route.Route     `json:"routes" yaml:"routes"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	APIKey  route.APIKey      `json:"api_key" yaml:"api-key"`
	CORS    bool              `json:"cors" yaml:"cors"`
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

func Load(p string, f *flags.Flags) (Config, error) {
	var c Config
	err := store.Load(p, &c)
	if err != nil {
		return Config{}, err
	}
	return c.merge(f), nil
}

func Store(p string, c Config) error {
	return store.Store(p, c)
}

func (c Config) merge(f *flags.Flags) Config {
	c.Debug = c.Debug || f.Debug
	c.Verbose = c.Verbose || f.Verbose

	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	for k, v := range f.Headers {
		c.Headers[k] = v
	}
	if v := f.APIKey; v.Valid() {
		c.APIKey = v
	}

	rts := make([]route.Route, 0, len(c.Routes)+len(f.Routes))
	for _, e := range c.Routes {
		e = e.WithHeaders(c.Headers)
		if c.APIKey.Valid() {
			e = e.WithAPIKey(c.APIKey)
		}
		rts = append(rts, e)
	}
	for _, e := range f.Routes {
		e = e.WithHeaders(c.Headers)
		if c.APIKey.Valid() {
			e = e.WithAPIKey(c.APIKey)
		}
		rts = append(rts, e)
	}
	c.Routes = rts

	return c
}
