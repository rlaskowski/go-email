package registries

import (
	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/grpc"
	"github.com/rlaskowski/go-email/queue"
)

type Registries struct {
	QueueFactory   *queue.QueueFactory
	Email          *email.Email
	GrpcEmailQueue *grpc.EmailQueue
}

func NewRegistries(queueFactory *queue.QueueFactory, email *email.Email) *Registries {
	r := &Registries{
		QueueFactory: queueFactory,
		Email:        email,
	}

	r.GrpcEmailQueue = grpc.NewEmailQueue(r.QueueFactory, r.Email)

	return r
}
