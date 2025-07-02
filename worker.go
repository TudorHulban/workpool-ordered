package workpoolordered

func (l *CList[T]) worker() {
	for node := range l.chWork {
		processedPayload, toDelete, errProcess := l.processor(node.Payload)
		if errProcess != nil {
			// TODO: log it locally
		}

		node.processedPayload = &processedPayload
		node.markedForDeletion.Store(toDelete)

		l.unprocessed.Add(-1)
	}
}
