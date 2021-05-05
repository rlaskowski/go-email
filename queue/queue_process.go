package queue

type QueueProcess interface {
	Start() error
	Stop() error
	Len() int
	Less(i, j int) bool
	Push(qstore interface{})
	Pop() interface{}
	Swap(i, j int)
}
