package config

import "flag"

func initFlags() {
	flag.IntVar(&HttpServerPort, "p", HttpServerPort, "Http server port")
	flag.BoolVar(&FileStore, "a", FileStore, "If true attached file is store on system disk")
	flag.StringVar(&FileStorePath, "f", FileStorePath, "Path where to store file before send")

	flag.Parse()
}
