package registries

import (
	"github.com/rlaskowski/go-email/queue"
)

type Registries struct {
	QueueFactory *queue.QueueFactory
}

func NewRegistries(QueueFactory *queue.QueueFactory) *Registries {
	return &Registries{QueueFactory}
}
