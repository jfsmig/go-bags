// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"golang.org/x/exp/constraints"
	"sort"
)

// SortedObj implements a sorted array of objects providing a PRIMARY KEY.
// The sorted storage allows O(log N) lookups, efficient sorted scans but
// average O(N log N) insertions.
// SortedObj proves useful for stable collection which are frequently accessed
// for paginated listings.
type SortedObj[PkType constraints.Ordered, T WithPK[PkType]] []T

type WithPK[PkType constraints.Ordered] interface {
	PK() PkType
}

func (s SortedObj[PkType, T]) Len() int { return len(s) }

func (s SortedObj[PkType, T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortedObj[PkType, T]) Less(i, j int) bool { return s[i].PK() < s[j].PK() }

func (s *SortedObj[PkType, T]) Add(a T) {
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

func (s SortedObj[PkType, T]) Slice(marker PkType, max uint32) []T {
	if max < minSliceSize {
		max = minSliceSize
	} else if max > maxSliceSize {
		max = maxSliceSize
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

func (s SortedObj[PkType, T]) Has(id PkType) bool { return s.GetIndex(id) >= 0 }

// Remove forwards the call to RemovePK with the primary key of the given
// element. The PK() method must manage the case of a `nil` receiver.
func (s *SortedObj[PkType, T]) Remove(a T) {
	s.RemovePK(a.PK())
}

// RemovePK identifies the position of the element with the given primary key
// and then removes it and restores the sorting of the set.
func (s *SortedObj[PkType, T]) RemovePK(pk PkType) {
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
