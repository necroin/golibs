package concurrent

import "sync"

type AtomicValue[T any] struct {
	value T
	mutex *sync.RWMutex
}

func NewAtomicValue[T any]() *AtomicValue[T] {
	return &AtomicValue[T]{
		mutex: &sync.RWMutex{},
	}
}

func (atomic *AtomicValue[T]) Get() T {
	atomic.mutex.RLock()
	defer atomic.mutex.RUnlock()
	return atomic.value
}

func (atomic *AtomicValue[T]) Set(value T) {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value = value
}
