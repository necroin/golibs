package concurrent

import "sync"

type AtomicBool struct {
	value bool
	mutex *sync.RWMutex
}

func NewAtomicBool() *AtomicBool {
	return &AtomicBool{
		mutex: &sync.RWMutex{},
	}
}

func (atomic *AtomicBool) Get() bool {
	atomic.mutex.RLock()
	defer atomic.mutex.RUnlock()
	return atomic.value
}

func (atomic *AtomicBool) Set(value bool) {
	atomic.mutex.Lock()
	defer atomic.mutex.Unlock()
	atomic.value = value
}

func (atomic *AtomicBool) Equal(other *AtomicBool) bool {
	selfValue := atomic.Get()
	otherValue := other.Get()
	return selfValue == otherValue
}

func (atomic *AtomicBool) NotEqual(other *AtomicBool) bool {
	return !atomic.Equal(other)
}
