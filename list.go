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

	dl.mutex.Lock()
	var pending []*Node[T]
	for node := dl.tail; node != nil; node = node.Prev {
		if node.ProcessedPayload == nil {
			pending = append(pending, node)
		}
	}
	dl.mutex.Unlock()

	if len(pending) == 0 {
		return nil
	}

	workerCount := dl.workerLimit
	if len(pending) < workerCount {
		workerCount = len(pending)
	}

	var wg sync.WaitGroup
	work := make(chan *Node[T], workerCount)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case node, ok := <-work:
					if !ok {
						return
					}

					result, err := dl.processor(node.Payload)
					if err == nil {
						dl.mutex.Lock()
						node.ProcessedPayload = &result
						dl.unprocessed--
						dl.mutex.Unlock()
					}
				}
			}
		}()
	}

	// Process from tail to head (oldest first)
	for i := 0; i < len(pending); i++ {
		select {
		case <-ctx.Done():
			close(work)
			wg.Wait()
			return ctx.Err()
		case work <- pending[i]:
		}
	}

	close(work)
	wg.Wait()
	return nil
}
