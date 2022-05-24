// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
	"testing"
)

func TestRaw_RemoveSorted(T *testing.T) {
	bag := SortedRaw[int]{0, 1, 2, 3}
	for _, v := range []int{-1, 4} {
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
	for _, v := range []int{0, 1, 2, 3} {
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
}

func TestRaw_RemoveUnsorted(T *testing.T) {
	bag := SortedRaw[int]{0, 1, 2, 3}
	for _, v := range []int{-1, 4} {
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
	for _, v := range []int{3, 2, 1, 0} {
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
}

func TestRaw_Lookup(T *testing.T) {
	bag := SortedRaw[int]{0, 1, 2, 3}
	for idx, v := range []int{0, 1, 2, 3} {
		if !bag.Has(v) {
			T.Fatal()
		}
		if x, ok := bag.Get(v); !ok {
			T.Fatal()
		} else if x != v {
			T.Fatal()
		}
		if idx != bag.GetIndex(v) {
			T.Fatal()
		}
	}
	for _, v := range []int{-1, -2, 5, 6} {
		if bag.Has(v) {
			T.Fatal()
		}
		if x, ok := bag.Get(v); ok {
			T.Fatal()
		} else if x != 0 {
			T.Fatal()
		}
		if -1 != bag.GetIndex(v) {
			T.Fatal()
		}
	}
}

func TestRaw_Slice(T *testing.T) {
	bag := SortedRaw[int]{0, 1, 2, 3}
	testSlice := func(marker int, max uint32, expectations ...int) {
		slice := bag.Slice(marker, max)
		if len(slice) != len(expectations) {
			T.Fatal()
		}
		for i, v := range slice {
			if v != expectations[i] {
				T.Fatal()
			}
		}
	}
	testSlice(bag[0], 2, bag[1], bag[2])
	testSlice(bag[0]-1, 2, bag[0], bag[1])
	testSlice(bag[len(bag)-1], 1)
	testSlice(bag[0]-1, minSliceSize-1, bag[:minSliceSize]...)
	testSlice(bag[0]-1, maxSliceSize+1, bag...)
}

func TestRaw_SliceLarge(T *testing.T) {
	const total = 3 * maxSliceSize
	var bag SortedRaw[int]
	for i := 0; i < total; i++ {
		bag.Add(i)
	}
	bag.Assert()
	slice := bag.Slice(bag[0], total)
	if len(slice) < minSliceSize || len(slice) > maxSliceSize {
		T.Fatal()
	}
}

func TestRaw_AppendUnsorted(T *testing.T) {
	var bag SortedRaw[int]
	bag.Append(3, 1, 0, 2)
	T.Log("bag", bag)
	bag.Assert()
	if bag.Len() != 4 {
		T.Fatalf("invalid length")
	}
}

func TestRaw_AppendSorted(T *testing.T) {
	var bag SortedRaw[int]
	bag.Append(0, 1, 2, 3)
	T.Log("bag", bag)
	bag.Assert()
	if bag.Len() != 4 {
		T.Fatal()
	}
}

func TestRaw_AddUnsorted(T *testing.T) {
	var bag SortedRaw[int]
	bag.Add(3)
	bag.Add(1)
	bag.Add(0)
	bag.Add(2)
	T.Log("bag", bag)
	bag.Assert()
}

func TestRaw_AddSorted(T *testing.T) {
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
	if s.Len() < 2 {
		return true
	}
	lastId := s[0]
	for _, a := range s[1:] {
		if lastId == a {
			return false
		}
		lastId = a
	}
	return true
}
