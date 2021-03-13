package route

import (
	"errors"
	"net/url"
	"strings"
)

var errMalformedRoute = errors.New("Malformed route")

type Route struct {
	Addr    string            `json:"addr" yaml:"addr"`
	URL     string            `json:"url" yaml:"url"`
	Headers map[string]string `json:"headers" yaml:"headers"`
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
