package concurrent

import (
	"fmt"
	"sync"

	"github.com/necroin/golibs/libs/container"
)

type ConcurrentMap[K comparable, V any] struct {
	data  *container.Map[K, V]
	mutex *sync.RWMutex
}

// Constructs a new container.
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data:  container.NewMap[K, V](),
		mutex: &sync.RWMutex{},
	}
}

// Inserts element into the container, replace if the container already contain an element with an equivalent key.
func (concurrentMap *ConcurrentMap[K, V]) Insert(key K, value V) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	concurrentMap.data.Insert(key, value)
}

// Finds an element with key equivalent to key.
func (concurrentMap *ConcurrentMap[K, V]) Find(key K) (V, bool) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return concurrentMap.data.Find(key)
}

// Removes specified element from the container.
func (concurrentMap *ConcurrentMap[K, V]) Erase(key K) (V, bool) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	return concurrentMap.data.Erase(key)
}

// Iterates over elements of the container with specified handler.
func (concurrentMap *ConcurrentMap[K, V]) Iterate(handler func(key K, value V)) {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	concurrentMap.data.Iterate(handler)
}

// Returns the number of elements in the container.
func (concurrentMap *ConcurrentMap[K, V]) Size() int {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return concurrentMap.data.Size()
}

// Checks if the container has no elements.
func (concurrentMap *ConcurrentMap[K, V]) IsEmpty() bool {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return concurrentMap.data.IsEmpty()
}

// Returns slice of map keys.
func (concurrentMap *ConcurrentMap[K, V]) Keys() []K {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return concurrentMap.data.Keys()
}

// Returns slice of map values.
func (concurrentMap *ConcurrentMap[K, V]) Values() []V {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return concurrentMap.data.Values()
}

// Adds a key/value pair to the container if the key does not already exist.
// Returns the new value, or the existing value if the key already exists.
func (concurrentMap *ConcurrentMap[K, V]) GetOrAddByFunc(key K, valueFactory func(key K) V) (V, bool) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	return concurrentMap.data.GetOrAddByFunc(key, valueFactory)
}

// Adds a key/value pair to the container if the key does not already exist.
// Returns the new value, or the existing value if the key already exists.
func (concurrentMap *ConcurrentMap[K, V]) GetOrAdd(key K, value V) (V, bool) {
	concurrentMap.mutex.Lock()
	defer concurrentMap.mutex.Unlock()
	return concurrentMap.data.GetOrAdd(key, value)
}

func (concurrentMap *ConcurrentMap[K, V]) String() string {
	concurrentMap.mutex.RLock()
	defer concurrentMap.mutex.RUnlock()
	return fmt.Sprintf("(len = %d) %v", concurrentMap.Size(), concurrentMap.data)
}
