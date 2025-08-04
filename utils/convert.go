package utils

func SliceToSlice[T any, S any](container []T, handler func(value T) S) []S {
	result := []S{}
	for _, value := range container {
		result = append(result, handler(value))
	}
	return result
}

func MapToMap[K comparable, V any, M comparable, N any](container map[K]V, handler func(key K, value V) (M, N)) map[M]N {
	result := map[M]N{}

	for key, value := range container {
		newKey, newValue := handler(key, value)
		result[newKey] = newValue
	}

	return result
}

func MapToSlice[K comparable, V any, S any](container map[K]V, handler func(key K, value V) S) []S {
	result := []S{}

	for key, value := range container {
		result = append(result, handler(key, value))
	}

	return result
}

func SliceToSliceMap[S any, K comparable, V any](container []S, handler func(record S) (K, V)) map[K][]V {
	result := map[K][]V{}

	if container == nil {
		return result
	}

	for _, record := range container {
		key, value := handler(record)
		mappedRecords, ok := result[key]
		if !ok {
			mappedRecords = []V{}
		}
		mappedRecords = append(mappedRecords, value)
		result[key] = mappedRecords
	}
	return result
}

func SliceToMap[S any, K comparable, V any](container []S, handler func(record S) (K, V)) map[K]V {
	result := map[K]V{}

	if container == nil {
		return result
	}

	for _, record := range container {
		key, value := handler(record)
		result[key] = value
	}

	return result
}
