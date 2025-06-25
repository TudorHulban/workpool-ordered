package workpoolordered

import (
	"context"
	"errors"
	"sync"
)

type Processor[T any] func(T) (T, error)

type Node[T any] struct {
	Payload          T
	ProcessedPayload *T // Use pointer to distinguish nil from zero value

	Prev *Node[T]
	Next *Node[T]
}

type DLinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]

	mutex sync.Mutex

	processor Processor[T]

	length      int
	workerLimit int
	unprocessed int
}

func NewDLinkedList[T any](processor Processor[T], workers int) *DLinkedList[T] {
	return &DLinkedList[T]{
		processor:   processor,
		workerLimit: workers,
	}
}

func (dl *DLinkedList[T]) Insert(payload T) {
	dl.mutex.Lock()
	defer dl.mutex.Unlock()

	dl.length++
	dl.unprocessed++

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

		dl.length--
		current = prev
	}

	return results
}

func (dl *DLinkedList[T]) Process(ctx context.Context) error {
	if dl.processor == nil {
		return errors.New("processor not set")
	}

	// Main continuous loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Single lock to check unprocessed AND collect pending work
			dl.mutex.Lock()

			if dl.unprocessed == 0 {
				dl.mutex.Unlock()
				return nil // All done, exit
			}

			// Find pending work while we have the lock
			var pending []*Node[T]
			for node := dl.tail; node != nil; node = node.Prev {
				if node.ProcessedPayload == nil {
					pending = append(pending, node)
				}
			}
			dl.mutex.Unlock()

			// Determine worker count for this iteration
			workerCount := dl.workerLimit
			if len(pending) < workerCount {
				workerCount = len(pending)
			}

			// Start workers for this batch
			var wg sync.WaitGroup

			for ix := 0; ix < workerCount; ix++ {
				wg.Add(1)

				go func(nodeIndex int) {
					defer wg.Done()

					if nodeIndex < len(pending) {
						node := pending[nodeIndex]

						result, err := dl.processor(node.Payload)
						if err == nil {
							dl.mutex.Lock()
							node.ProcessedPayload = &result
							dl.unprocessed--
							dl.mutex.Unlock()
						}
					}
				}(ix)
			}

			wg.Wait()
		}
	}
}
