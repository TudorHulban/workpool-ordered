package workpoolordered

import (
	"errors"
	"sync"
	"sync/atomic"
)

type CList[T any] struct {
	head *Node[T]
	tail *Node[T]

	sync.Mutex

	processor Processor[T]
	chWork    chan *Node[T]
	chStop    chan struct{}

	unprocessed atomic.Int64
	length      atomic.Int64
}

type ParamsCList[T any] struct {
	Processor       Processor[T]
	Workers         int
	WaitToStartWork bool //introduced for tests.
}

func NewCList[T any](params *ParamsCList[T]) (*CList[T], error) {
	if params.Processor == nil {
		return nil,
			errors.New("processor not set")
	}

	result := CList[T]{
		processor: params.Processor,
		chWork:    make(chan *Node[T]),
		chStop:    make(chan struct{}),
	}

	for range params.Workers {
		go result.worker()
	}

	if !params.WaitToStartWork {
		go result.doWork()
	}

	return &result,
		nil
}

func (l *CList[T]) Close() {
	l.chStop <- struct{}{}
}
