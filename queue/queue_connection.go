package queue

type QueueConnection interface {
	Start() error
	Stop() error
	Publish(message interface{}) error
	Subscribe() error
}
