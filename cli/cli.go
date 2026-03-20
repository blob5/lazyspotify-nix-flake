package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/dubeyKartikay/lazyspotify/core/auth"
)

func Run(args []string){ 
	switch args[0] {
		case "auth":
			authHandler(args)
		case "test":
      testHandler(args)
	}
}

func authHandler(args []string) {
	if(len(args) > 1){
		printUsage()
		return;
	}
	ctx,cancel := context.WithTimeoutCause(context.Background(), 30*time.Second, fmt.Errorf("auth: timeout"))
	defer cancel()
	_, err := auth.New().Authenticate(ctx)
  if err != nil {
    fmt.Println("Error authenticating:", err)
    return
  }
  fmt.Println("Authenticated with Spotify")
}

func testHandler(args []string) {
	ctx,cancel := context.WithTimeoutCause(context.Background(), 30*time.Second, fmt.Errorf("helloworld: timeout"))
	defer cancel()
	client,err := auth.New().GetClient(ctx)
  if err != nil {
    return
  }
	fmt.Println(client.GetAlbum(context.Background(),"0kzl3HWoYqLTBaFJ3DjpqT"))
}

func printUsage(){
	fmt.Println("Usage: lazyspotify <command>")
	fmt.Println("Commands:")
	fmt.Println("  auth    Authenticate with Spotify")
}
