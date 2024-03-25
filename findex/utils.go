package findex

func StringSliceBytes(value []string) int64 {
	result := int64(0)
	for _, raw := range value {
		result += int64(len([]rune(raw)))
	}
	return result
}
