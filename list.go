package workpoolordered

import (
	"sync"
	"sync/atomic"
)

type DLinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]

	mutex sync.Mutex

	processor Processor[T]

	unprocessed atomic.Int64
	length      atomic.Int64
	workerLimit int
}

func NewDLinkedList[T any](processor Processor[T], workers int) *DLinkedList[T] {
	return &DLinkedList[T]{
		processor:   processor,
		workerLimit: workers,
	}
}

func (dl *DLinkedList[T]) Insert(payload T) {
	dl.length.Add(1)
	dl.unprocessed.Add(1)

	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	node := &Node[T]{
		Payload: payload,
	}

	if dl.head == nil {
		dl.head, dl.tail = node, node
		return
	}

	node.Next = dl.head
	dl.head.Prev = node
	dl.head = node
}

// Read returns all processed payloads up to the first unprocessed one
func (dl *DLinkedList[T]) Read() []T {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	var results []T
	current := dl.tail

	for current != nil {
		if current.ProcessedPayload == nil {
			break
		}
		results = append(results, *current.ProcessedPayload)
		prev := current.Prev

		// Remove processed node
		if current.Prev != nil {
			current.Prev.Next = current.Next
		} else {
			dl.head = current.Next
		}

		if current.Next != nil {
			current.Next.Prev = current.Prev
		} else {
			dl.tail = current.Prev
		}

		dl.length.Add(-1)
		current = prev
	}

	return results
}
