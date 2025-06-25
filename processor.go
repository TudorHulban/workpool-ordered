package workpoolordered

type Processor[T any] func(T) (T, error)
