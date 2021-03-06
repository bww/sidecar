package proxy

import (
	"log"
	"net/http"
	"time"

	"sidecar/route"
)

type Config struct {
	Routes  []route.Route
	CORS    bool
	Verbose bool
	Debug   bool
}

type Proxy struct {
	Config
	svcs []*http.Server
}

func NewWithConfig(conf Config) (*Proxy, error) {
	var svcs []*http.Server
	for _, e := range conf.Routes {
		h, err := newHandler(conf, e)
		if err != nil {
			return nil, err
		}
		s := &http.Server{
			Addr:           e.Addr,
			Handler:        h,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		svcs = append(svcs, s)
	}
	return &Proxy{conf, svcs}, nil
}

func (p *Proxy) ListenAndServe() error {
	errs := make(chan error)
	for _, e := range p.svcs {
		if p.Debug {
			log.Printf("Starting %s → %s", e.Addr, e.Handler.(Handler).Describe())
		} else if p.Verbose {
			log.Printf("Starting %s → %v", e.Addr, e.Handler)
		}
		go func(s *http.Server) {
			errs <- s.ListenAndServe()
		}(e)
	}
	return <-errs
}
