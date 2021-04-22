package queue

type QueueConnection interface {
	Start() error
	Stop() error
	Publish(Subject QueueSubject, message ...interface{}) error
	Subscribe(Subject QueueSubject) ([]QueueStore, error)
}
