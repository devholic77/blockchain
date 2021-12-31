package cli

import (
	"flag"

	"github.com/devholic77/duckcoin/explorer"
	"github.com/devholic77/duckcoin/rest"
)

func Start() {
	mode := flag.String("mode", "rest", "server mode")
	port := flag.Int("port", 4000, "server port")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	}
}
