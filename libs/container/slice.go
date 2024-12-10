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
func (slice *Slice[V]) Insert(index uint, value V) error {
	size := slice.Size()
	if int(index) >= size {
		return fmt.Errorf("index out of range")
	}

	slice.data[index] = value
	return nil
}

// Appends the given elements value to the end of the container.
func (slice *Slice[V]) Append(values ...V) {
	slice.data = append(slice.data, values...)
}

// Returns the element at specified location index, with bounds checking.
// If index is not within the range of the container, an error is returned.
func (slice *Slice[V]) At(index uint) (V, error) {
	size := slice.Size()
	if int(index) >= size {
		return *new(V), fmt.Errorf("index out of range")
	}

	return slice.data[index], nil
}

// Erases the specified element from the container.
func (slice *Slice[V]) Erase(index uint) error {
	size := slice.Size()

	if int(index) >= size {
		return fmt.Errorf("index out of range")
	}

	if int(index) == size-1 {
		slice.data = slice.data[0:index]
	} else {
		slice.data = append(slice.data[0:index], slice.data[index+1:size]...)
	}

	return nil
}

// Returns the number of elements in the container.
func (slice *Slice[V]) Size() int {
	return len(slice.data)
}

// Checks if the container has no elements.
func (slice *Slice[V]) IsEmpty() bool {
	return slice.Size() == 0
}

// Returns the first element in the container.
// Calling front on an empty container causes undefined behavior.
func (slice *Slice[V]) Front() V {
	if slice.Size() == 0 {
		return *new(V)
	}
	return slice.data[0]
}

// Returns the last element in the container.
// Calling back on an empty container causes undefined behavior.
func (slice *Slice[V]) Back() V {
	size := slice.Size()
	if size == 0 {
		return *new(V)
	}
	return slice.data[size-1]
}

// Returns the element at specified location index, with bounds checking.
// Erases the specified element from the container.
// If index is not within the range of the container, an error is returned.
func (slice *Slice[V]) PopAt(index uint) (V, error) {
	result, err := slice.At(index)
	if err != nil {
		return result, err
	}

	if err := slice.Erase(index); err != nil {
		return result, err
	}

	return result, nil
}

// Returns an iterator to the first element of the container.
func (slice *Slice[V]) Begin() *SliceIterator[V] {
	return &SliceIterator[V]{
		data:  slice,
		index: 0,
	}
}

// Returns an iterator to the element following the last element of the container.
func (slice *Slice[V]) End() *SliceIterator[V] {
	return &SliceIterator[V]{
		data:  slice,
		index: uint(slice.Size()),
	}
}

func (slice *Slice[V]) String() string {
	return fmt.Sprintf("(len = %d) %v", slice.Size(), slice.data)
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
