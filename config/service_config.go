package config

import "time"

type ServiceConfig struct {
	//Path to store temporary message
	FileStorePath string

	//Is message have to store to temporary file
	FileStore bool

	//Default HTTP server port
	HttpServerPort int

	//Maximum duration for reading request
	HttpServerReadTimeout time.Duration

	//Maximum duration before timing out
	//writes of the response
	HttpServerWriteTimeout time.Duration

	//Maximum number of bytes the
	//server will read parsing the request header's
	HttpMaxHeaderSize int

	//Default GRPC server port
	GrpcListenPort int

	//Duration when queue will be refreshing
	QueueRefreshTime time.Duration
}

const (
	ExecutableName  = "email"
	EmailConfigFile = "config.yaml"
	FilePermissions = 0700
	FileCopyBuff    = 1024 * 1024
)

var (
	DefaultServiceConfig = ServiceConfig{
		FileStorePath:          GetWorkingDirectory(),
		HttpServerPort:         8181,
		HttpServerReadTimeout:  30 * time.Second,
		HttpServerWriteTimeout: 30 * time.Second,
		HttpMaxHeaderSize:      1024 * 4,
		GrpcListenPort:         9090,
		QueueRefreshTime:       5 * time.Second,
	}
)
