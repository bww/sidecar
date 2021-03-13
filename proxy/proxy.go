package proxy

import (
	"net/http"
	"time"

	"sidecar/route"
)

type Proxy struct {
	svcs []*http.Server
}

func NewWithRoutes(r []route.Route) (*Proxy, error) {
	var svcs []*http.Server
	for _, e := range r {
		h, err := e.Handler()
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
	return &Proxy{svcs}, nil
}

func (p *Proxy) ListenAndServe() error {
	errs := make(chan error)
	for _, e := range p.svcs {
		go func(s *http.Server) {
			errs <- s.ListenAndServe()
		}(e)
	}
	return <-errs
}
