package utils

import (
	"math/rand"
)

func GetRandomFrom[T any](values ...T) T {
	return values[rand.Intn(len(values))]
}
