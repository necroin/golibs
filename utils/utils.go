package utils

import (
	"fmt"
	"os"
	"strings"
)

func PointerOf[T any](value T) *T {
	return &value
}

func InstantiateSliceElement[T any](value *[]T) *T {
	return new(T)
}

func CleanTag(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return parts[0]
}

func MapCopy[K comparable, V any](value map[K]V) map[K]V {
	result := map[K]V{}
	for k, v := range value {
		result[k] = v
	}
	return result
}

func MapKeys[K comparable, V any](container map[K]V) []K {
	result := []K{}
	for key := range container {
		result = append(result, key)
	}
	return result
}

func MapValues[K comparable, V any](container map[K]V) []V {
	result := []V{}
	for _, value := range container {
		result = append(result, value)
	}
	return result
}

func SaveToFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed open file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed write data to file: %w", err)
	}

	return nil
}
