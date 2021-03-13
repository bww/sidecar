package main

import (
	"fmt"
	"os"

	"sidecar/config"
	"sidecar/flags"
	"sidecar/proxy"
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
	flag := flags.New(args[0])
	err := flag.Parse(args[1:])
	if err != nil {
		return err
	}

	var conf config.Config
	if p := flag.Config; p != "" {
		conf, err = config.Load(p, flag)
		if err != nil {
			return err
		}
	}

	pxy, err := proxy.NewWithConfig(proxy.Config{
		Routes:  conf.Routes,
		CORS:    conf.CORS,
		Verbose: conf.Verbose,
		Debug:   conf.Debug,
	})
	if err != nil {
		return err
	}

	panic(pxy.ListenAndServe())
	return nil
}
