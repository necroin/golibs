package utils

func PointerOf[T any](value T) *T {
	return &value
}

func InstantiateSliceElement[T any](value *[]T) *T {
	return new(T)
}
