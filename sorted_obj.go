// Copyright (c) 2018-2023 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
)

// SortedObj implements a sorted array of objects providing a PRIMARY KEY.
// The sorted storage allows O(log N) lookups, efficient sorted scans but
// average O(N log N) insertions.
// SortedObj proves useful for stable collection which are frequently accessed
// for paginated listings.
type SortedObj[PkType Ordered, T WithPK[PkType]] []T

type WithPK[PkType Ordered] interface {
	PK() PkType
}

// Len implements a method of the sort.Interface
func (s SortedObj[PkType, T]) Len() int { return len(s) }

// Swap implements a method of the sort.Interface
func (s SortedObj[PkType, T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements a method of the sort.Interface
func (s SortedObj[PkType, T]) Less(i, j int) bool { return s[i].PK() < s[j].PK() }

// Add introduces a new item in the sorted array, regardless the presence of another item with the same PRIMARY KEY
// and preserves the ordering of the array.
func (s *SortedObj[PkType, T]) Add(a T) {
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

// Append introduces several items in the sorted array, regardless the presence of other items with the same PRIMARY KEY
// and preserves the ordering of the array.
func (s *SortedObj[PkType, T]) Append(a ...T) {
	*s = append(*s, a...)
	sort.Sort(s)
}

func (s SortedObj[PkType, T]) Slice(marker PkType, max uint32) []T {
	if max < MinSliceSize {
		max = MinSliceSize
	} else if max > MaxSliceSize {
		max = MaxSliceSize
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].PK() > marker
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

// GetIndex returns -1 if no item of the array has the given PRIMARY KEY, or the position of the first
// element with that PRIMARY KEY
func (s SortedObj[PkType, T]) GetIndex(id PkType) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].PK() >= id
	})
	if i < len(s) && s[i].PK() == id {
		return i
	}
	return -1
}

func (s SortedObj[PkType, T]) Get(id PkType) (out T, ok bool) {
	idx := s.GetIndex(id)
	if idx >= 0 {
		return s[idx], true
	}
	return out, false
}

// Has tests for the presence of an item in the set, given the private key of the item
func (s SortedObj[PkType, T]) Has(id PkType) bool { return s.GetIndex(id) >= 0 }

// Remove identifies the position of the element with the given PRIMARY KEY
// and then removes it and restores the sorting of the set.
func (s *SortedObj[PkType, T]) Remove(pk PkType) {
	idx := s.GetIndex(pk)
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
