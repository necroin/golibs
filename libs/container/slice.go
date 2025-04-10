package container

import (
	"fmt"
)

type Slice[V any] struct {
	data []V
}

type SliceIterator[V any] struct {
	data  *Slice[V]
	index uint
}

// Constructs a new container.
func NewSlice[V any]() *Slice[V] {
	return &Slice[V]{
		data: []V{},
	}
}

// Inserts element at the specified location in the container.
func (container *Slice[V]) Insert(index uint, value V) error {
	size := container.Size()
	if int(index) >= size {
		return fmt.Errorf("index out of range")
	}

	container.data[index] = value
	return nil
}

// Appends the given elements value to the end of the container.
func (container *Slice[V]) Append(values ...V) {
	container.data = append(container.data, values...)
}

// Returns the element at specified location index, with bounds checking.
// If index is not within the range of the container, an error is returned.
func (container *Slice[V]) At(index uint) (V, error) {
	size := container.Size()
	if int(index) >= size {
		return *new(V), fmt.Errorf("index out of range")
	}

	return container.data[index], nil
}

// Erases the specified element from the container.
func (container *Slice[V]) Erase(index uint) error {
	size := container.Size()

	if int(index) >= size {
		return fmt.Errorf("index out of range")
	}

	if int(index) == size-1 {
		container.data = container.data[0:index]
	} else {
		container.data = append(container.data[0:index], container.data[index+1:size]...)
	}

	return nil
}

// Returns the number of elements in the container.
func (container *Slice[V]) Size() int {
	return len(container.data)
}

// Checks if the container has no elements.
func (container *Slice[V]) IsEmpty() bool {
	return container.Size() == 0
}

// Returns the first element in the container.
// Calling front on an empty container causes undefined behavior.
func (container *Slice[V]) Front() V {
	if container.Size() == 0 {
		return *new(V)
	}
	return container.data[0]
}

// Returns the last element in the container.
// Calling back on an empty container causes undefined behavior.
func (container *Slice[V]) Back() V {
	size := container.Size()
	if size == 0 {
		return *new(V)
	}
	return container.data[size-1]
}

// Returns the element at specified location index, with bounds checking.
// Erases the specified element from the container.
// If index is not within the range of the container, an error is returned.
func (container *Slice[V]) PopAt(index uint) (V, error) {
	result, err := container.At(index)
	if err != nil {
		return result, err
	}

	if err := container.Erase(index); err != nil {
		return result, err
	}

	return result, nil
}

// Returns an iterator to the first element of the container.
func (container *Slice[V]) Begin() *SliceIterator[V] {
	return &SliceIterator[V]{
		data:  container,
		index: 0,
	}
}

// Returns an iterator to the element following the last element of the container.
func (container *Slice[V]) End() *SliceIterator[V] {
	return &SliceIterator[V]{
		data:  container,
		index: uint(container.Size()),
	}
}

func (container *Slice[V]) String() string {
	return fmt.Sprintf("(len = %d) %v", container.Size(), container.data)
}

func (iterator *SliceIterator[V]) Next() *SliceIterator[V] {
	iterator.index += 1
	return iterator
}

func (iterator *SliceIterator[V]) Get() (V, error) {
	return iterator.data.At(iterator.index)
}

func (iterator *SliceIterator[V]) Pos() uint {
	return iterator.index
}

func (iterator *SliceIterator[V]) Set(value V) error {
	return iterator.data.Insert(iterator.index, value)
}

func (iterator *SliceIterator[V]) Equal(other *SliceIterator[V]) bool {
	selfPos := iterator.Pos()
	otherPos := other.Pos()
	return selfPos == otherPos
}
