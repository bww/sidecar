package route

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var errMalformedRoute = errors.New("Malformed route")

type APIKey struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
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

func (r Route) Handler() (Handler, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return Handler{}, err
	}

	h := make(http.Header)
	for k, v := range r.Headers {
		h.Set(k, v)
	}
	if r.APIKey.Username != "" {
		h.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(r.APIKey.Username+":"+r.APIKey.Password)))
	}

	v := Handler{
		proxy:   &httputil.ReverseProxy{},
		url:     u,
		headers: h,
	}

	v.proxy.Director = v.director

	return v, nil
}

type Handler struct {
	proxy   *httputil.ReverseProxy
	url     *url.URL
	headers http.Header
}

func (h Handler) director(req *http.Request) {
	target := h.url

	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)

	if target.RawQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
	}

	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "") // explicitly disable User-Agent so it's not set to default value
	}
	for k, v := range h.headers {
		for _, e := range v {
			req.Header.Set(k, e)
		}
	}
}

func (h Handler) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	h.proxy.ServeHTTP(rsp, req)
}
