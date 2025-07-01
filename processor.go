package workpoolordered

// Processor returns in bool if the node is marked for deletion.
type Processor[T any] func(T) (T, bool, error)
