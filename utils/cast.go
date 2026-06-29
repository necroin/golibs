package utils

import (
	"strconv"
)

func PointerOf[T any](value T) *T {
	return &value
}

func GetPointerValueOrDefault[T any](value *T) T {
	if value == nil {
		return *new(T)
	}
	return *value
}

func CastToInt(value *string) (*int, error) {
	if value == nil {
		return nil, nil
	}

	if *value == "" {
		return PointerOf(0), nil
	}

	result, castErr := strconv.ParseInt(*value, 10, 32)

	if castErr != nil {
		if floatResult, err := CastToFloat32(value); floatResult != nil && err == nil {
			return PointerOf(int(*floatResult)), nil
		}
	}

	return PointerOf(int(result)), castErr
}

func CastToInt32(value *string) (*int32, error) {
	if value == nil {
		return nil, nil
	}

	if *value == "" {
		return PointerOf(int32(0)), nil
	}

	result, castErr := strconv.ParseInt(*value, 10, 32)

	if castErr != nil {
		if floatResult, err := CastToFloat32(value); floatResult != nil && err == nil {
			return PointerOf(int32(*floatResult)), nil
		}
	}

	return PointerOf(int32(result)), castErr
}

func CastToInt64(value *string) (*int64, error) {
	if value == nil {
		return nil, nil
	}

	if *value == "" {
		return PointerOf(int64(0)), nil
	}

	result, castErr := strconv.ParseInt(*value, 10, 64)

	if castErr != nil {
		if floatResult, err := CastToFloat64(value); floatResult != nil && err == nil {
			return PointerOf(int64(*floatResult)), nil
		}
	}

	return PointerOf(result), castErr
}

func CastToUInt(value *string) (*uint, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return PointerOf(uint(0)), nil
	}
	result, err := strconv.ParseInt(*value, 10, 32)
	return PointerOf(uint(result)), err
}

func CastToUInt32(value *string) (*uint32, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return PointerOf(uint32(0)), nil
	}
	result, err := strconv.ParseInt(*value, 10, 32)
	return PointerOf(uint32(result)), err
}

func CastToUInt64(value *string) (*uint64, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return PointerOf(uint64(0)), nil
	}
	result, err := strconv.ParseInt(*value, 10, 64)
	return PointerOf(uint64(result)), err
}

func CastToFloat32(value *string) (*float32, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return PointerOf(float32(0)), nil
	}
	result, err := strconv.ParseFloat(*value, 32)
	return PointerOf(float32(result)), err
}

func CastToFloat64(value *string) (*float64, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return PointerOf(float64(0)), nil
	}
	result, err := strconv.ParseFloat(*value, 64)
	return PointerOf(result), err
}

func CastToBool(value *string) (*bool, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return PointerOf(false), nil
	}
	boolResult, err := strconv.ParseBool(*value)
	return PointerOf(boolResult), err
}
