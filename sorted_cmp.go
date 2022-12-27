// Copyright (c) 2018-2023 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
)

// SortedCmp implements a sorted array of comparable objects.
// The sorted storage allows O(log N) lookups, efficient sorted scans but
// average O(N log N) insertions.
// SortedCmp proves useful for stable collection which are frequently accessed
// for paginated listings.
type SortedCmp[T WithCompare[T]] []T

type WithCompare[T any] interface {
	Compare(other T) int
}

// Len implements a method of the sort.Interface
func (s SortedCmp[T]) Len() int { return len(s) }

// Swap implements a method of the sort.Interface
func (s SortedCmp[T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements a method of the sort.Interface
func (s SortedCmp[T]) Less(i, j int) bool { return s[i].Compare(s[j]) < 0 }

// Add introduces a new item in the sorted array, regardless the presence of the same item,
// and preserves the ordering of the array.
func (s *SortedCmp[T]) Add(a T) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		// Only sort the array if the last 2 element are not sorted: in other words,
		// adding the new biggest element maintains the ordering
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

// Append introduces several items in the sorted array, regardless the presence of identical  items
// and preserves the ordering of the array.
func (s *SortedCmp[T]) Append(a ...T) {
	*s = append(*s, a...)
	sort.Sort(s)
}

func (s SortedCmp[T]) Slice(marker T, max uint32) []T {
	if max < MinSliceSize {
		max = MinSliceSize
	} else if max > MaxSliceSize {
		max = MaxSliceSize
	}
	start := sort.Search(len(s), func(i int) bool {
		return marker.Compare(s[i]) < 0
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

// GetIndex returns the position of the first items that matches (Compare returns 0) to the given other item,
// or -1 in case of no match.
func (s SortedCmp[T]) GetIndex(id T) int {
	i := sort.Search(len(s), func(i int) bool {
		return id.Compare(s[i]) <= 0
	})
	if i < len(s) && id.Compare(s[i]) == 0 {
		return i
	}
	return -1
}

func (s SortedCmp[T]) Get(id T) (out T, ok bool) {
	idx := s.GetIndex(id)
	if idx >= 0 {
		return s[idx], true
	}
	return out, false
}

// Has tests for the presence of an item in the set, given a copy of the item
func (s SortedCmp[T]) Has(id T) bool { return s.GetIndex(id) >= 0 }

// Remove forwards the call to RemovePK with the primary key of the given
// element.
func (s *SortedCmp[T]) Remove(a T) {
	idx := s.GetIndex(a)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}
