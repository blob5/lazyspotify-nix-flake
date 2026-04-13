package cli

import (
	"fmt"
	"os"

	"github.com/dubeyKartikay/lazyspotify/buildinfo"
)

func Run(args []string) {
	switch args[0] {
	case "auth":
		authHandler(args)
	case "play":
		playHandler(args)
	case "version":
		versionHandler(args)
	default:
		printUsage()
	}
}

func authHandler(args []string) {
	if len(args) > 1 {
		printUsage()
		return
	}

}

func playHandler(args []string) {
	if len(args) != 2 {
		printUsage()
		return
	}
}

func versionHandler(args []string) {
	if len(args) != 1 {
		printUsage()
		return
	}

	_ = buildinfo.PrintVersion(os.Stdout)
}

func printUsage() {
	fmt.Println("Usage: lazyspotify <command>")
	fmt.Println("Flags:")
	fmt.Println("  --version Print build metadata")
	fmt.Println("Commands:")
	fmt.Println("  auth    Authenticate with Spotify")
	fmt.Println("  play    Play a hardcoded Spotify track")
	fmt.Println("  version Print build metadata")
}
