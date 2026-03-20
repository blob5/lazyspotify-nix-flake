package main
import ("flag"
				"fmt"
				"github.com/dubeyKartikay/lazyspotify/cli"
)
func main() {
	flag.Parse()
	switch {
		case flag.NArg() > 0:
      cli.Run(flag.Args())
    default:
      fmt.Println("Running the TUI")
	}
}
