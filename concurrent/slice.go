package concurrent

import (
	"fmt"
	"sync"
)

type ConcurrentSlice[V any] struct {
	data  []V
	mutex *sync.RWMutex
}

type ConcurrentSliceIterator[V any] struct {
	data  *ConcurrentSlice[V]
	index uint
	mutex *sync.RWMutex
}

// Constructs a new container.
func NewConcurrentSlice[V any]() *ConcurrentSlice[V] {
	return &ConcurrentSlice[V]{
		data:  []V{},
		mutex: &sync.RWMutex{},
	}
}

// Inserts element at the specified location in the container.
func (concurrentSlice *ConcurrentSlice[V]) Insert(index uint, value V) error {
	if int(index) >= len(concurrentSlice.data) {
		return fmt.Errorf("[ConcurrentSlice] [Error] index out of range")
	}

	concurrentSlice.mutex.Lock()
	defer concurrentSlice.mutex.Unlock()

	concurrentSlice.data[index] = value

	return nil
}

// Appends the given element value to the end of the container.
func (concurrentSlice *ConcurrentSlice[V]) Append(values ...V) {
	concurrentSlice.mutex.Lock()
	defer concurrentSlice.mutex.Unlock()
	concurrentSlice.data = append(concurrentSlice.data, values...)
}

// Returns the element at specified location index, with bounds checking.
// If index is not within the range of the container, an error is returned.
func (concurrentSlice *ConcurrentSlice[V]) At(index uint) (V, error) {
	if int(index) >= len(concurrentSlice.data) {
		return *new(V), fmt.Errorf("[ConcurrentSlice] [Error] index out of range")
	}

	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()

	return concurrentSlice.data[index], nil
}

// Erases the specified elements from the container.
func (concurrentSlice *ConcurrentSlice[V]) Erase(index uint) error {
	size := concurrentSlice.Size()

	if int(index) >= size {
		return fmt.Errorf("[ConcurrentSlice] [Error] index out of range")
	}

	concurrentSlice.mutex.Lock()
	defer concurrentSlice.mutex.Unlock()

	if int(index) == size-1 {
		concurrentSlice.data = concurrentSlice.data[0:index]
	} else {
		concurrentSlice.data = append(concurrentSlice.data[0:index], concurrentSlice.data[index+1:size]...)
	}

	return nil
}

// Returns the number of elements in the container.
func (concurrentSlice *ConcurrentSlice[V]) Size() int {
	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()
	return len(concurrentSlice.data)
}

// Checks if the container has no elements.
func (concurrentSlice *ConcurrentSlice[V]) IsEmpty() bool {
	return concurrentSlice.Size() == 0
}

// Returns the first element in the container.
// Calling front on an empty container causes undefined behavior.
func (concurrentSlice *ConcurrentSlice[V]) Front() V {
	if concurrentSlice.Size() == 0 {
		return *new(V)
	}

	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()

	return concurrentSlice.data[0]
}

// Returns the last element in the container.
// Calling back on an empty container causes undefined behavior.
func (concurrentSlice *ConcurrentSlice[V]) Back() V {
	size := concurrentSlice.Size()
	if size == 0 {
		return *new(V)
	}

	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()

	return concurrentSlice.data[size-1]
}

// Returns an iterator to the first element of the container.
func (concurrentSlice *ConcurrentSlice[V]) Begin() *ConcurrentSliceIterator[V] {
	return &ConcurrentSliceIterator[V]{
		data:  concurrentSlice,
		index: 0,
		mutex: &sync.RWMutex{},
	}
}

// Returns an iterator to the element following the last element of the container.
func (concurrentSlice *ConcurrentSlice[V]) End() *ConcurrentSliceIterator[V] {
	return &ConcurrentSliceIterator[V]{
		data:  concurrentSlice,
		index: uint(concurrentSlice.Size()),
		mutex: &sync.RWMutex{},
	}
}

func (iterator *ConcurrentSliceIterator[V]) Next() *ConcurrentSliceIterator[V] {
	iterator.mutex.Lock()
	defer iterator.mutex.Unlock()
	iterator.index += 1
	return iterator
}

func (iterator *ConcurrentSliceIterator[V]) Get() (V, error) {
	iterator.mutex.RLock()
	defer iterator.mutex.RUnlock()
	return iterator.data.At(iterator.index)
}

func (iterator *ConcurrentSliceIterator[V]) Set(value V) error {
	iterator.mutex.Lock()
	defer iterator.mutex.Unlock()
	return iterator.data.Insert(iterator.index, value)
}

func (iterator *ConcurrentSliceIterator[V]) Equal(other *ConcurrentSliceIterator[V]) bool {
	iterator.mutex.RLock()
	other.mutex.RLock()
	defer iterator.mutex.RUnlock()
	defer other.mutex.RUnlock()
	return iterator.index == other.index
}
