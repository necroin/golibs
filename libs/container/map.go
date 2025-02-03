package container

import "fmt"

type Map[K comparable, V any] struct {
	data map[K]V
}

// Constructs a new container.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		data: map[K]V{},
	}
}

// Inserts element into the container, replace if the container already contain an element with an equivalent key.
func (container *Map[K, V]) Insert(key K, value V) {
	container.data[key] = value
}

// Finds an element with key equivalent to key.
func (container *Map[K, V]) Find(key K) (V, bool) {
	result, ok := container.data[key]
	if !ok {
		result = *new(V)
	}
	return result, ok
}

// Removes specified element from the container.
func (container *Map[K, V]) Erase(key K) (V, bool) {
	result, ok := container.data[key]
	if !ok {
		result = *new(V)
	}
	delete(container.data, key)
	return result, ok
}

// Iterates over elements of the container with specified handler.
func (container *Map[K, V]) Iterate(handler func(key K, value V)) {
	for key, value := range container.data {
		handler(key, value)
	}
}

// Returns the number of elements in the container.
func (container *Map[K, V]) Size() int {
	return len(container.data)
}

// Checks if the container has no elements.
func (container *Map[K, V]) IsEmpty() bool {
	return container.Size() == 0
}

// Returns slice of map keys.
func (container *Map[K, V]) Keys() []K {
	result := []K{}
	container.Iterate(func(key K, value V) {
		result = append(result, key)
	})
	return result
}

// Returns slice of map values.
func (container *Map[K, V]) Values() []V {
	result := []V{}
	container.Iterate(func(key K, value V) {
		result = append(result, value)
	})
	return result
}

// Adds a key/value pair to the container if the key does not already exist.
// Returns the new value, or the existing value if the key already exists.
func (container *Map[K, V]) GetOrAddByFunc(key K, valueFactory func(key K) V) (V, bool) {
	result, ok := container.data[key]
	if !ok {
		result = valueFactory(key)
		container.data[key] = result
	}
	return result, ok
}

// Adds a key/value pair to the container if the key does not already exist.
// Returns the new value, or the existing value if the key already exists.
func (container *Map[K, V]) GetOrAdd(key K, value V) (V, bool) {
	return container.GetOrAddByFunc(key, func(key K) V { return value })
}

func (container *Map[K, V]) String() string {
	return fmt.Sprintf("(len = %d) %v", container.Size(), container.data)
}
