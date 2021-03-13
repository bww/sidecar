package main

import (
	"fmt"
	"os"

	"sidecar/config"
	"sidecar/proxy"
	"sidecar/route"
)

var ( // set at compile time via the linker
	version = "v0.0.0"
	githash = "000000"
)

func main() {
	if err := app(os.Args); err != nil {
		fmt.Println("***", err)
		os.Exit(1)
	}
}

func app(args []string) error {
	flag := newFlags(args[0])
	err := flag.Parse(args[1:])
	if err != nil {
		return err
	}

	var conf config.Config
	if p := flag.Config; p != "" {
		conf, err = config.Load(p)
		if err != nil {
			return err
		}
	}

	conf.Debug = conf.Debug || flag.Debug
	conf.Verbose = conf.Verbose || flag.Verbose

	hdrs := make(map[string]string)
	for k, v := range conf.Headers {
		hdrs[k] = v
	}
	for k, v := range flag.Headers {
		hdrs[k] = v
	}

	var apikey route.APIKey
	if conf.APIKey.Valid() {
		apikey = conf.APIKey
	}
	if flag.APIKey.Valid() {
		apikey = conf.APIKey
	}

	rts := make([]route.Route, 0, len(conf.Routes)+len(flag.Routes))
	for _, e := range conf.Routes {
		rts = append(rts, e.WithHeaders(hdrs).WithAPIKey(apikey))
	}
	for _, e := range flag.Routes {
		rts = append(rts, e.WithHeaders(hdrs).WithAPIKey(apikey))
	}
	if len(rts) < 1 {
		return fmt.Errorf("No routes defined; nothing to do")
	}

	pxy, err := proxy.NewWithRoutes(proxy.Config{
		Verbose: conf.Verbose,
		Debug:   conf.Debug,
	}, rts)
	if err != nil {
		return err
	}

	panic(pxy.ListenAndServe())
	return nil
}
