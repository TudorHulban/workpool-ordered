package workpoolordered

// Read returns all processed payloads up to the first unprocessed one.
// Deletes all nodes marked for deletion.
func (l *CList[T]) Read() []T {
	l.Lock()
	defer l.Unlock()

	var results []T
	current := l.tail

	for current != nil {
		if current.IsMarkedForDeletion() {
			l.delete(current)

			continue
		}

		if current.processedPayload == nil {
			current = current.prev

			continue
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
