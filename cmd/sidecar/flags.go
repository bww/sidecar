package main

import (
	"flag"
	"fmt"
	"strings"

	"sidecar/route"
)

type Flags struct {
	*flag.FlagSet

	Debug   bool
	Verbose bool
	Quiet   bool
	Config  string
	Routes  []route.Route
	Headers map[string]string

	Values struct {
		Debug   *bool
		Verbose *bool
		Quiet   *bool
		Config  *string
		Routes  flagList
		Headers flagList
	}
}

func newFlags(cmd string) *Flags {
	f := &Flags{
		FlagSet: flag.NewFlagSet(cmd, flag.ExitOnError),
	}

	f.Values.Debug = f.Bool("debug", false, "Be extremely verbose.")
	f.Values.Verbose = f.Bool("verbose", false, "Be more verbose.")
	f.Values.Config = f.String("config", "", "The configuration file to use.")

	f.Var(&f.Values.Routes, "route", "Define a route that maps a local port to a remote base URI, specified as '<port>:<uri>'. Provide -route repeatedly to set many routes.")
	f.Var(&f.Values.Headers, "header", "Define a header to be set on every proxied request; specified as '<name>: <value>'. Provide -header repeatedly to set many headers.")
	return f
}

func (f *Flags) Parse(args []string) error {
	f.FlagSet.Parse(args)

	f.Debug = *f.Values.Debug
	f.Verbose = *f.Values.Verbose
	f.Config = *f.Values.Config

	var routes []route.Route
	for _, e := range f.Values.Routes {
		r, err := route.Parse(e)
		if err != nil {
			return err
		}
		routes = append(routes, r)
	}

	headers := make(map[string]string)
	for _, e := range f.Values.Headers {
		if x := strings.Index(e, ":"); x > 0 {
			headers[strings.TrimSpace(e[:x])] = strings.TrimSpace(e[x+1:])
		}
	}

	f.Routes = routes
	f.Headers = headers
	return nil
}

type flagList []string

func (s *flagList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *flagList) String() string {
	return fmt.Sprintf("%+v", *s)
}
