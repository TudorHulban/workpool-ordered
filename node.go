package workpoolordered

type Node[T any] struct {
	Payload          T
	processedPayload *T // Use pointer to distinguish nil from zero value

	prev *Node[T]
	next *Node[T]

	markedForDeletion bool
}
