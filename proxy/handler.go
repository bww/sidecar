package proxy

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"sidecar/route"
)

type Handler struct {
	Config
	proxy   *httputil.ReverseProxy
	url     *url.URL
	headers http.Header
}

func newHandler(conf Config, rte route.Route) (Handler, error) {
	u, err := url.Parse(rte.URL)
	if err != nil {
		return Handler{}, err
	}

	h := make(http.Header)
	for k, v := range rte.Headers {
		h.Set(k, v)
	}
	if rte.APIKey.Username != "" {
		h.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(rte.APIKey.Username+":"+rte.APIKey.Password)))
	}

	v := Handler{
		Config:  conf,
		url:     u,
		headers: h,
	}
	v.proxy = &httputil.ReverseProxy{
		Director: v.director,
	}

	return v, nil
}

func (h Handler) director(req *http.Request) {
	src := fmt.Sprintf("%s%v", req.Host, req.URL)
	target := h.url

	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)

	if target.RawQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
	}

	req.Host = target.Host

	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "") // explicitly disable User-Agent so it's not set to default value
	}
	for k, v := range h.headers {
		for _, e := range v {
			req.Header.Set(k, e)
		}
	}

	if h.Verbose || h.Debug {
		log.Printf("%v → %v", src, req.URL)
	}
}

func (h Handler) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	h.proxy.ServeHTTP(rsp, req)
}

func (h Handler) String() string {
	return h.url.String()
}

func (h Handler) Describe() string {
	b := &strings.Builder{}
	b.WriteString(h.url.String())
	for k, v := range h.headers {
		b.WriteString("\n\t")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(strings.Join(v, ", "))
	}
	return b.String()
}
