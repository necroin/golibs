package concurrent

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	data           map[K]V
	mutex          *sync.RWMutex
	complexOpMutex *sync.Mutex
}

// Constructs a new container.
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data:           map[K]V{},
		mutex:          &sync.RWMutex{},
		complexOpMutex: &sync.Mutex{},
	}
}

// Inserts element into the container, replace if the container already contain an element with an equivalent key.
func (concurrentMap *ConcurrentMap[K, V]) Insert(key K, value V) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	concurrentMap.data[key] = value
}

// Finds an element with key equivalent to key.
func (concurrentMap *ConcurrentMap[K, V]) Find(key K) (V, bool) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	result, ok := concurrentMap.data[key]
	if !ok {
		result = *new(V)
	}
	return result, ok
}

// Removes specified element from the container.
func (concurrentMap *ConcurrentMap[K, V]) Erase(key K) (V, bool) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	result, ok := concurrentMap.data[key]
	if !ok {
		result = *new(V)
	}
	delete(concurrentMap.data, key)
	return result, ok
}

// Iterates over elements of the container with specified handler.
func (concurrentMap *ConcurrentMap[K, V]) Iterate(handler func(key K, value V)) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()

	for key, value := range concurrentMap.data {
		handler(key, value)
	}
}

// Returns the number of elements in the container.
func (concurrentMap *ConcurrentMap[K, V]) Size() int {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return len(concurrentMap.data)
}

// Checks if the container has no elements.
func (concurrentMap *ConcurrentMap[K, V]) IsEmpty() bool {
	return concurrentMap.Size() == 0
}

// Returns slice of map keys.
func (concurrentMap *ConcurrentMap[K, V]) Keys() []K {
	result := []K{}
	concurrentMap.Iterate(func(key K, value V) {
		result = append(result, key)
	})
	return result
}

// Returns slice of map values.
func (concurrentMap *ConcurrentMap[K, V]) Values() []V {
	result := []V{}
	concurrentMap.Iterate(func(key K, value V) {
		result = append(result, value)
	})
	return result
}

// Executes complex operation on map with given handler.
func (concurrentMap *ConcurrentMap[K, V]) ComplexOperation(handler func() error) error {
	concurrentMap.complexOpMutex.Lock()
	defer concurrentMap.complexOpMutex.Unlock()
	return handler()
}
