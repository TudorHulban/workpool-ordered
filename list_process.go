package workpoolordered

// doWork is blocking so should be run in own goroutine.
func (l *CList[T]) doWork() error {
	for {
		select {
		case <-l.chStop:
			return nil

		default:
			l.Lock()

			if l.unprocessed.Load() == 0 {
				l.Unlock()

				return nil // All done
			}

			var pending []*Node[T]

			for node := l.tail; node != nil; node = node.prev {
				if node.processedPayload == nil {
					pending = append(pending, node)
				}
			}

			l.Unlock()

			if len(pending) == 0 {
				continue
			}

			// distribute found work
			for _, node := range pending {
				l.chWork <- node
			}
		}
	}
}
