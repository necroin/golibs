package concurrent

import (
	"sync"

	"github.com/necroin/golibs/utils"
)

type AtomicNumber[T utils.Number] struct {
	value T
	mutex *sync.RWMutex
}

func NewAtomicNumber[T utils.Number]() *AtomicNumber[T] {
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

func (atomic *AtomicNumber[T]) Add(value T) T {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value += value
	return atomic.value
}

func (atomic *AtomicNumber[T]) Sub(value T) T {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value -= value
	return atomic.value
}

func (atomic *AtomicNumber[T]) Inc() T {
	return atomic.Add(1)

}

func (atomic *AtomicNumber[T]) Dec() T {
	return atomic.Sub(1)
}
