package main

import (
	"fmt"
	"os"

	"sidecar/config"
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
	flag := parseFlags(args[0])

	conf, err := config.Find()
	if err != nil {
		return err
	}

	_ = flag
	_ = conf
	return nil
}
