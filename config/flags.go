package config

import "flag"

func initFlags() {
	flag.IntVar(&HttpServerPort, "p", HttpServerPort, "Http server port")

	flag.Parse()
}
