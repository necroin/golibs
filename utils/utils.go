package utils

func PointerOf[T any](value T) *T {
	return &value
}

func InstantiateSliceElement[T any](value *[]T) *T {
	return new(T)
}

func MapCopy[K comparable, V any](value map[K]V) map[K]V {
	result := map[K]V{}
	for k, v := range value {
		result[k] = v
	}
	return result
}
