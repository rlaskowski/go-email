package config

import (
	"time"
)

var (
	ExecutableName         = "email"
	HttpServerPort         = 8080
	HttpServerReadTimeout  = 30 * time.Second
	HttpServerWriteTimeout = 30 * time.Second
	HttpMaxHeaderSize      = 1024 * 4
)
