package utils

import (
	"fmt"
	"math"
)

type Number interface {
	float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func SafeDivide[T Number](v1, v2 T) T {
	if v2 == 0 {
		return 0
	}
	return v1 / v2
}

type Vector2D struct {
	X, Y float64
}

func (vector *Vector2D) Distance() float64 {
	return math.Sqrt(vector.X*vector.X + vector.Y*vector.Y)
}

func (vector *Vector2D) String() string {
	return fmt.Sprintf("{ X: %v, Y: %v }", vector.X, vector.Y)
}
