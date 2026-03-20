package cli

import (
	"fmt"
	"github.com/dubeyKartikay/lazyspotify/core/auth"
  "context"
)

func Run(args []string){ 
	switch args[0] {
		case "auth":
			authHandler(args)
		case "test":
      helloworldHandler(args)
	}
}

func authHandler(args []string) {
	if(len(args) > 1){
		printUsage()
		return;
	}
	_, err := auth.New().Authenticate()
  if err != nil {
    fmt.Println("Error authenticating:", err)
    return
  }
  fmt.Println("Authenticated with Spotify")
}

func helloworldHandler(args []string) {
	client,err := auth.New().GetClient()
  if err != nil {
    return
  }
  fmt.Println("Hello World")
	fmt.Println(client.GetAlbum(context.Background(),"0kzl3HWoYqLTBaFJ3DjpqT"))
}

func printUsage(){
	fmt.Println("Usage: lazyspotify <command>")
	fmt.Println("Commands:")
	fmt.Println("  auth    Authenticate with Spotify")
}
