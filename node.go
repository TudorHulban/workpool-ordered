package workpoolordered

type Node[T any] struct {
	Payload          T
	ProcessedPayload *T // Use pointer to distinguish nil from zero value

	Prev *Node[T]
	Next *Node[T]
}
