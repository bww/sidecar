package route

import (
	"errors"
	"strings"
)

var errMalformedRoute = errors.New("Malformed route")

type APIKey struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (a APIKey) Valid() bool {
	return a.Username != ""
}

type Route struct {
	Addr    string            `json:"addr" yaml:"addr"`
	URL     string            `json:"url" yaml:"url"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	APIKey  APIKey            `json:"api_key" yaml:"api_key"`
}

func Parse(s string) (Route, error) {
	var addr, base string
	if x := strings.Index(s, "="); x > 0 {
		addr, base = s[:x], s[x+1:]
	} else {
		return Route{}, errMalformedRoute
	}
	return Route{
		Addr: addr,
		URL:  base,
	}, nil
}

func (r Route) WithAPIKey(k APIKey) Route {
	return Route{
		Addr:    r.Addr,
		URL:     r.URL,
		Headers: r.Headers,
		APIKey:  k,
	}
}

func (r Route) WithHeaders(h map[string]string) Route {
	d := make(map[string]string)
	for k, v := range r.Headers {
		d[k] = v
	}
	for k, v := range h {
		d[k] = v
	}
	return Route{
		Addr:    r.Addr,
		URL:     r.URL,
		Headers: d,
		APIKey:  r.APIKey,
	}
}
