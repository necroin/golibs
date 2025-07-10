package utils

import (
	"math/rand"
	"time"
)

func GetRandomFrom[T any](values ...T) T {
	return values[rand.Intn(len(values))]
}

func NewRandomGenerator() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return generator
}
