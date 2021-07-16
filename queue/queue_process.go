package queue

type QueueProcess interface {
	Len() int
	Less(i, j int) bool
	Push(qstore interface{})
	Pop() interface{}
	Swap(i, j int)
}
