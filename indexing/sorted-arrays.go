package indexing

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type item[K constraints.Ordered, V any] struct {
	key   K
	value V
}

type SortedArray[K constraints.Ordered, V any] struct {
	items []item[K, V]
}

func NewSortedArray[K constraints.Ordered, V any](items []item[K, V]) *SortedArray[K, V] {
	sort.Slice(items, func(i, j int) bool {
		return items[i].key < items[j].key
	})
	return &SortedArray[K, V]{items: items}
}
