package workpoolordered

import (
	"errors"
	"sync"
	"sync/atomic"
)

type CList[T any] struct {
	head *Node[T]
	tail *Node[T]

	sync.Mutex

	processor Processor[T]
	chWork    chan *Node[T]
	chStop    chan struct{}

	unprocessed atomic.Int64
	length      atomic.Int64
}

type ParamsCList[T any] struct {
	Processor       Processor[T]
	Workers         int
	WaitToStartWork bool //introduced for tests.
}

func NewCList[T any](params *ParamsCList[T]) (*CList[T], error) {
	if params.Processor == nil {
		return nil,
			errors.New("processor not set")
	}

	result := CList[T]{
		processor: params.Processor,
		chWork:    make(chan *Node[T]),
		chStop:    make(chan struct{}),
	}

	for range params.Workers {
		go result.worker()
	}

	if !params.WaitToStartWork {
		go result.process()
	}

	return &result,
		nil
}

func (l *CList[T]) Insert(payload T) {
	l.length.Add(1)
	l.unprocessed.Add(1)

	l.Lock()
	defer l.Unlock()

	node := &Node[T]{
		Payload: payload,
	}

	if l.head == nil {
		l.head, l.tail = node, node
		return
	}

	node.next = l.head
	l.head.prev = node
	l.head = node
}

// Read returns all processed payloads up to the first unprocessed one
func (l *CList[T]) Read() []T {
	l.Lock()
	defer l.Unlock()

	var results []T
	current := l.tail

	for current != nil {
		if current.processedPayload == nil {
			break
		}
		results = append(results, *current.processedPayload)
		prev := current.prev

		// Remove processed node
		if current.prev != nil {
			current.prev.next = current.next
		} else {
			l.head = current.next
		}

		if current.next != nil {
			current.next.prev = current.prev
		} else {
			l.tail = current.prev
		}

		l.length.Add(-1)
		current = prev
	}

	return results
}
