package registries

import (
	"github.com/rlaskowski/go-email/grpc"
	"github.com/rlaskowski/go-email/queue"
)

type Registries struct {
	queueBox       *queue.QueueBox
	GrpcEmailQueue *grpc.EmailService
}

func NewRegistries() *Registries {
	r := &Registries{
		queueBox: new(queue.QueueBox),
	}

	r.GrpcEmailQueue = grpc.NewEmailService(r.queueBox)

	return r
}

func (r *Registries) Start() error {
	if err := r.queueBox.Start(); err != nil {
		return err
	}

	return nil
}

func (r *Registries) Stop() error {
	if err := r.queueBox.Stop(); err != nil {
		return err
	}

	return nil
}
