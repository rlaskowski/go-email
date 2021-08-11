package registry

import (
	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/grpc"
	"github.com/rlaskowski/go-email/queue"
)

type Registry interface {
	EmailService() *grpc.EmailService

	QueueBox() *queue.QueueBox

	ServiceConfig() config.ServiceConfig
}
