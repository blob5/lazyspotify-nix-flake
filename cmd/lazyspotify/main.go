package main

import (
	"flag"
	"os"

	"github.com/dubeyKartikay/lazyspotify/buildinfo"
	"github.com/dubeyKartikay/lazyspotify/cli"
	ui "github.com/dubeyKartikay/lazyspotify/ui/v1"
)

var versionFlag bool

func init() {
	flag.BoolVar(&versionFlag, "version", false, "print build metadata")
}

func main() {
	flag.Parse()
	if versionFlag {
		_ = buildinfo.PrintVersion(os.Stdout)
		return
	}

	switch {
	case flag.NArg() > 0:
		cli.Run(flag.Args())
	default:
		ui.RunTui()
	}
}
