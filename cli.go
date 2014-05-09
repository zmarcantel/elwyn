package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var Opts Options
var Verbosity int

type Options struct {
	// helpers
	Verbosity int // save as an int, not slice

	// options
	Verbose []bool `short:"v"   long:"verbose"        description:"Show verbose log information. Supports -v[vvv] syntax."`
	Quiet   bool   `short:"q"   long:"quiet"          description:"Do not output any text to STDOUT."`
	LogPath string `short:"l"   long:"log"            description:"Path to folder where logs should be saved" default:"./logs"`
	Port    int    `short:"p"   long:"port"           description:"Port to listen for HTTP traffic" default:"8765"`
}

func cliInit() error {
	if _, help := flags.Parse(&Opts); help != nil {
		os.Exit(1)
	}

	// handle verbosity
	Opts.Verbosity = len(Opts.Verbose)

	return nil
}
