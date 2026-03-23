package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/dubeyKartikay/lazyspotify/core/auth"
	"github.com/dubeyKartikay/lazyspotify/core/player"
	"github.com/dubeyKartikay/lazyspotify/core/utils"
)

func Run(args []string) {
	switch args[0] {
	case "auth":
		authHandler(args)
	case "play":
		playHandler(args)
  default:
    printUsage()
	}
}

func authHandler(args []string) {
	if len(args) > 1 {
		printUsage()
		return
	}
	ctx, cancel := context.WithTimeoutCause(context.Background(), 30*time.Second, fmt.Errorf("auth: timeout"))
	defer cancel()
	_, err := auth.New().ReAuthenticate(ctx)
	if err != nil {
		fmt.Println("Error authenticating:", err)
		return
	}
	fmt.Println("Authenticated with Spotify")
}

func playHandler(args []string) {
	if len(args) != 2{
		printUsage()
		return
	}
	uri,err := utils.SpotifyURLToURI(args[1])
  if err != nil {
    fmt.Println("Error parsing spotify url:", err)
    return
  }
	if err := player.PlayTrack(context.Background(),uri); err != nil {
		fmt.Println("Error playing hardcoded track:", err)
	}
}

func printUsage() {
	fmt.Println("Usage: lazyspotify <command>")
	fmt.Println("Commands:")
	fmt.Println("  auth    Authenticate with Spotify")
	fmt.Println("  play    Play a hardcoded Spotify track")
}
