// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
)

// SortedCmp implements a sorted array of objects providing a PRIMARY KEY.
// The sorted storage allows O(log N) lookups, efficient sorted scans but
// average O(N log N) insertions.
// SortedCmp proves useful for stable collection which are frequently accessed
// for paginated listings.
type SortedCmp[T WithCompare[T]] []T

type WithCompare[T any] interface {
	Compare(other T) int
}

func (s SortedCmp[T]) Len() int { return len(s) }

func (s SortedCmp[T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortedCmp[T]) Less(i, j int) bool { return s[i].Compare(s[j]) < 0 }

func (s *SortedCmp[T]) Add(a T) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
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

func (s *SortedCmp[T]) Append(a ...T) {
	*s = append(*s, a...)
	sort.Sort(s)
}

func (s SortedCmp[T]) Slice(marker T, max uint32) []T {
	if max < minSliceSize {
		max = minSliceSize
	} else if max > maxSliceSize {
		max = maxSliceSize
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
