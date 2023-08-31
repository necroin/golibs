package concurrent

import (
	"sync"
)

type Number interface {
	float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type AtomicNumber[T Number] struct {
	value T
	mutex *sync.RWMutex
}

func NewAtomicNumber[T Number]() *AtomicNumber[T] {
	return &AtomicNumber[T]{
		mutex: &sync.RWMutex{},
	}
}

func (atomic *AtomicNumber[T]) Get() T {
	atomic.mutex.RLock()
	defer atomic.mutex.RUnlock()
	return atomic.value
}

func (atomic *AtomicNumber[T]) Set(value T) {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value = value
}

func (atomic *AtomicNumber[T]) Add(value T) {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value += value
}

func (atomic *AtomicNumber[T]) Sub(value T) {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value -= value
}
