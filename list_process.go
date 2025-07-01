package workpoolordered

import (
	"errors"
	"sync"
)

func (dl *DLinkedList[T]) Process() error {
	if dl.processor == nil {
		return errors.New("processor not set")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			dl.mutex.Lock()

			if dl.unprocessed.Load() == 0 {
				dl.mutex.Unlock()
				return nil // All done
			}

			var pending []*Node[T]
			for node := dl.tail; node != nil; node = node.Prev {
				if node.ProcessedPayload == nil {
					pending = append(pending, node)
				}
			}

			dl.mutex.Unlock()

			if len(pending) == 0 {
				continue
			}

			// Limit to max allowed workers
			workerCount := dl.workerLimit
			if len(pending) < workerCount {
				workerCount = len(pending)
			}

			var wg sync.WaitGroup
			wg.Add(workerCount)

			for ix := range workerCount {
				node := pending[ix] // safe slice access

				go func(n *Node[T]) {
					defer wg.Done()

					result, errProcess := dl.processor(n.Payload)
					if errProcess == nil {
						dl.mutex.Lock()
						n.ProcessedPayload = &result
						dl.unprocessed.Add(-1)
						dl.mutex.Unlock()
					}
				}(node)
			}

			wg.Wait()
		}
	}
}
