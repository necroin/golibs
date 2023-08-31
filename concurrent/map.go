package concurrent

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	data  map[K]V
	mutex *sync.RWMutex
}

// Constructs a new container.
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data:  map[K]V{},
		mutex: &sync.RWMutex{},
	}
}

// Inserts element into the container, if the container doesn't already contain an element with an equivalent key.
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
func (concurrentMap *ConcurrentMap[K, V]) Iterate(handler func(K, V)) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()

	for key, value := range concurrentMap.data {
		handler(key, value)
	}
}
