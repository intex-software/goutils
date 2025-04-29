package maps

import (
	"cmp"
	"slices"
)

func SortedMapKeys[S ~map[E]V, E cmp.Ordered, V any](input S) (keys []E) {
	keys = make([]E, 0, len(input))
	for i := range input {
		keys = append(keys, i)
	}
	slices.Sort(keys)
	return
}
