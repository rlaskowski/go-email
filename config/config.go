package config

import (
	"time"
)

const (
	ExecutableName  = "email"
	EmailConfigFile = "config.yaml"
	FilePermissions = 0700
	FileCopyBuff    = 1024 * 1024
)

var (
	FileStorePath          = GetWorkingDirectory()
	FileStore              = false
	HttpServerPort         = 8181
	HttpServerReadTimeout  = 30 * time.Second
	HttpServerWriteTimeout = 30 * time.Second
	HttpMaxHeaderSize      = 1024 * 4
	GrpcListenPort         = 9090
	QueueRefreshTime       = 5 * time.Second
)
