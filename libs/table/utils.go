package table

import (
	"slices"

	"github.com/necroin/golibs/utils"
)

func StringsListWitdth(values ...string) int {
	result := 0
	for _, value := range values {
		result += len(value)
	}

	return result
}

func TableRowWidth(values ...string) int {
	return StringsListWitdth(values...) + len(values) - 1
}

func TableRowsMaxWidth(rows ...[]string) int {
	return slices.Max(utils.SliceToSlice(rows, func(values []string) int { return TableRowWidth(values...) }))
}
