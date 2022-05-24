// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
	"testing"
)

type CmpInt int

func (i0 CmpInt) Compare(i1 CmpInt) int {
	return int(i0) - int(i1)
}

func TestCmp_Unsorted(T *testing.T) {
	var bag SortedCmp[CmpInt]
	bag.Add(3)
	bag.Add(1)
	bag.Add(0)
	bag.Add(2)
	T.Log("bag", bag)
	bag.Assert()
}

func TestCmp_Sorted(T *testing.T) {
	var bag SortedCmp[CmpInt]
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
func (s SortedCmp[T]) Assert() {
	if err := s.Check(); err != nil {
		panic(err)
	}
}

// Check validates the ordering and the unicity of the elements in the array
func (s SortedCmp[T]) Check() error {
	if !sort.IsSorted(s) {
		return ErrUnsorted
	}
	if !s.areItemsUnique() {
		return ErrDuplicates
	}
	return nil
}

// areItemsUnique validates the unicity of the elements in the array
func (s SortedCmp[T]) areItemsUnique() bool {
	lastId := s[0]
	for _, a := range s[1:] {
		if lastId.Compare(a) == 0 {
			return false
		}
		lastId = a
	}
	return true
}
