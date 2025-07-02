package workpoolordered

func (l *CList[T]) Insert(payload T) {
	l.length.Add(1)
	l.unprocessed.Add(1)

	l.Lock()
	defer l.Unlock()

	// move node to heap.
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
