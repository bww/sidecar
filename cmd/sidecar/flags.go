package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	*flag.FlagSet

	Debug   bool
	Verbose bool
	Quiet   bool
	Project string
	Headers map[string]string
	Routes  []string

	Values struct {
		Debug   *bool
		Verbose *bool
		Quiet   *bool
		Headers flagList
		Routes  flagList
	}
}

func parseFlags(cmd string) *Flags {
	f := &Flags{
		FlagSet: flag.NewFlagSet(cmd, flag.ExitOnError),
	}

	f.Values.Debug = f.Bool("debug", false, "Be extremely verbose.")
	f.Values.Verbose = f.Bool("verbose", false, "Be more verbose.")
	f.Values.Quiet = f.Bool("quiet", false, "Only print minimal output.")

	f.Var(&f.Values.Routes, "route", "Define a route that maps a local port to a remote base URI, specified as '<port>:<uri>'. Provide -route repeatedly to set many routes.")
	f.Var(&f.Values.Headers, "header", "Define a header to be set on every proxied request; specified as '<name>: <value>'. Provide -header repeatedly to set many headers.")
	return f
}

func (f *Flags) Parse(args []string) {
	f.FlagSet.Parse(args)
	f.Debug = *f.Values.Debug
	f.Verbose = *f.Values.Verbose
	f.Quiet = *f.Values.Quiet
}

type flagList []string

func (s *flagList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *flagList) String() string {
	return fmt.Sprintf("%+v", *s)
}
