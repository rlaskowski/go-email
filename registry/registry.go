package registry

import (
	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/grpc"
	"github.com/rlaskowski/go-email/queue"
	"github.com/rlaskowski/go-email/rest"
)

type Registry interface {
	EmailService() *grpc.EmailService
	EmailRestService() *rest.EmailService
	QueueBox() *queue.QueueBox
	ServiceConfig() config.ServiceConfig
}
