package concurrent

import (
	"fmt"
	"sync"

	"github.com/necroin/golibs/libs/container"
)

type ConcurrentSlice[V any] struct {
	data  *container.Slice[V]
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
		data:  container.NewSlice[V](),
		mutex: &sync.RWMutex{},
	}
}

// Inserts element at the specified location in the container.
func (concurrentSlice *ConcurrentSlice[V]) Insert(index uint, value V) error {
	concurrentSlice.mutex.Lock()
	defer concurrentSlice.mutex.Unlock()
	return concurrentSlice.data.Insert(index, value)
}

// Appends the given elements value to the end of the container.
func (concurrentSlice *ConcurrentSlice[V]) Append(values ...V) {
	concurrentSlice.mutex.Lock()
	defer concurrentSlice.mutex.Unlock()
	concurrentSlice.data.Append(values...)
}

// Returns the element at specified location index, with bounds checking.
// If index is not within the range of the container, an error is returned.
func (concurrentSlice *ConcurrentSlice[V]) At(index uint) (V, error) {
	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()
	return concurrentSlice.data.At(index)
}

// Erases the specified element from the container.
func (concurrentSlice *ConcurrentSlice[V]) Erase(index uint) error {
	concurrentSlice.mutex.Lock()
	defer concurrentSlice.mutex.Unlock()
	return concurrentSlice.data.Erase(index)
}

// Returns the number of elements in the container.
func (concurrentSlice *ConcurrentSlice[V]) Size() int {
	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()
	return concurrentSlice.data.Size()
}

// Checks if the container has no elements.
func (concurrentSlice *ConcurrentSlice[V]) IsEmpty() bool {
	return concurrentSlice.Size() == 0
}

// Returns the first element in the container.
// Calling front on an empty container causes undefined behavior.
func (concurrentSlice *ConcurrentSlice[V]) Front() V {
	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()
	return concurrentSlice.data.Front()
}

// Returns the last element in the container.
// Calling back on an empty container causes undefined behavior.
func (concurrentSlice *ConcurrentSlice[V]) Back() V {
	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()
	return concurrentSlice.data.Back()
}

// Returns the element at specified location index, with bounds checking.
// Erases the specified element from the container.
// If index is not within the range of the container, an error is returned.
func (concurrentSlice *ConcurrentSlice[V]) PopAt(index uint) (V, error) {
	concurrentSlice.mutex.RLock()
	defer concurrentSlice.mutex.RUnlock()
	return concurrentSlice.data.PopAt(index)
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

func (concurrentSlice *ConcurrentSlice[V]) String() string {
	return fmt.Sprintf("(len = %d) %v", concurrentSlice.Size(), concurrentSlice.data)
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

func (iterator *ConcurrentSliceIterator[V]) Pos() uint {
	iterator.mutex.RLock()
	defer iterator.mutex.RUnlock()
	return iterator.index
}

func (iterator *ConcurrentSliceIterator[V]) Set(value V) error {
	iterator.mutex.Lock()
	defer iterator.mutex.Unlock()
	return iterator.data.Insert(iterator.index, value)
}

func (iterator *ConcurrentSliceIterator[V]) Equal(other *ConcurrentSliceIterator[V]) bool {
	selfPos := iterator.Pos()
	otherPos := other.Pos()
	return selfPos == otherPos
}
