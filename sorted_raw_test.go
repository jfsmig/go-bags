// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
	"testing"
)

func TestRaw_Unsorted(T *testing.T) {
	var bag SortedRaw[int]
	bag.Add(3)
	bag.Add(1)
	bag.Add(0)
	bag.Add(2)
	T.Log("bag", bag)
	bag.Assert()
}

func TestRaw_Sorted(T *testing.T) {
	var bag SortedRaw[int]
	bag.Add(0)
	bag.Add(1)
	bag.Add(2)
	bag.Add(3)
	T.Log("bag", bag)
	bag.Assert()

	s := bag.Slice(1, 2)
	T.Log("slice", s)
	if len(s) != 2 {
		panic("bad length")
	} else if s[0] != 2 && s[1] != 3 {
		panic("bad content")
	}
}

// Assert panics if Check returns an error
func (s SortedRaw[int]) Assert() {
	if err := s.Check(); err != nil {
		panic(err)
	}
}

// Check validates the ordering and the unicity of the elements in the array
func (s SortedRaw[int]) Check() error {
	if !sort.IsSorted(s) {
		return ErrUnsorted
	}
	if !s.areItemsUnique() {
		return ErrDuplicates
	}
	return nil
}

// areItemsUnique validates the unicity of the elements in the array
func (s SortedRaw[int]) areItemsUnique() bool {
	lastId := s[0]
	for _, a := range s[1:] {
		if lastId == a {
			return false
		}
		lastId = a
	}
	return true
}
