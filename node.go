package workpoolordered

import "sync/atomic"

type Node[T any] struct {
	Payload          T
	processedPayload *T // Use pointer to distinguish nil from zero value

	prev *Node[T]
	next *Node[T]

	markedForDeletion atomic.Bool
}

func (n *Node[T]) MarkForDeletion() {
	n.markedForDeletion.Store(true)
}

func (n *Node[T]) IsProcessed() bool {
	return n.processedPayload != nil
}

func (n *Node[T]) IsMarkedForDeletion() bool {
	return n.markedForDeletion.Load()
}
