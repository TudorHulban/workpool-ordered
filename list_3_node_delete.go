package workpoolordered

// Delete removes a node from the list and returns true if successful.
// If the node is not found, it returns false.
func (l *CList[T]) Delete(node *Node[T]) bool {
	if node == nil {
		return false
	}

	l.Lock()
	defer l.Unlock()

	// Special case: empty list
	if l.head == nil {
		return false
	}

	// Special case: deleting the head node
	if node == l.head {
		if l.head == l.tail { // Only one node in list
			l.head = nil
			l.tail = nil
		} else {
			l.head = l.head.next
			l.tail.next = l.head // Maintain circularity
		}

		if !node.IsProcessed() {
			l.unprocessed.Add(-1)
		}
		l.length.Add(-1)

		return true
	}

	// Find the node and its predecessor
	prev := l.head
	current := l.head.next

	for current != l.head && current != node {
		prev = current
		current = current.next
	}

	// Node not found
	if current == l.head {
		return false
	}

	// Update links
	prev.next = current.next

	// If we are deleting the tail, update tail pointer
	if current == l.tail {
		l.tail = prev
	}

	if !node.IsProcessed() {
		l.unprocessed.Add(-1)
	}

	l.length.Add(-1)

	return true
}

// delete internal, to be used in Read.
func (l *CList[T]) delete(node *Node[T]) bool {
	// Special case: deleting the head node
	if node == l.head {
		if l.head == l.tail { // Only one node in list
			l.head = nil
			l.tail = nil
		} else {
			l.head = l.head.next
			l.tail.next = l.head // Maintain circularity
		}

		if !node.IsProcessed() {
			l.unprocessed.Add(-1)
		}
		l.length.Add(-1)

		return true
	}

	// Find the node and its predecessor
	prev := l.head
	current := l.head.next

	for current != l.head && current != node {
		prev = current
		current = current.next
	}

	// Node not found
	if current == l.head {
		return false
	}

	// Update links
	prev.next = current.next

	// If we are deleting the tail, update tail pointer
	if current == l.tail {
		l.tail = prev
	}

	if !node.IsProcessed() {
		l.unprocessed.Add(-1)
	}

	l.length.Add(-1)

	return true
}
