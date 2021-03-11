package route

import (
	"errors"
	"net/url"
	"strings"
)

var errMalformedRoute = errors.New("Malformed route")

type Route struct {
	Addr string
	URL  *url.URL
}

func Parse(s string) (Route, error) {
	var addr, base string
	if x := strings.Index(s, "="); x > 0 {
		addr, base = s[:x], s[x+1:]
	} else {
		return Route{}, errMalformedRoute
	}
	u, err := url.Parse(base)
	if err != nil {
		return Route{}, err
	}
	return Route{
		Addr: addr,
		URL:  u,
	}, nil
}
